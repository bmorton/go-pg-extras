package pgextras

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// Format controls how query results are rendered.
type Format int

const (
	FormatTable Format = iota // ASCII bordered table
	FormatJSON                // JSON array
	FormatCSV                 // CSV with header row
)

// FormatResults writes results in the given format.
// results must be a slice of result structs.
func FormatResults[T any](w io.Writer, results []T, format Format, title string) error {
	switch format {
	case FormatTable:
		return formatTable(w, results, title)
	case FormatJSON:
		return formatJSON(w, results)
	case FormatCSV:
		return formatCSV(w, results)
	default:
		return fmt.Errorf("unsupported format: %d", format)
	}
}

func formatTable[T any](w io.Writer, results []T, title string) error {
	headers, rows := extractHeadersAndRows(results)

	table := tablewriter.NewTable(w,
		tablewriter.WithRowAutoWrap(tw.WrapNone),
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)

	headerArgs := make([]any, len(headers))
	for i, h := range headers {
		headerArgs[i] = h
	}
	table.Header(headerArgs...)

	if title != "" {
		table.Caption(tw.Caption{Text: title})
	}

	for _, row := range rows {
		if err := table.Append(row); err != nil {
			return err
		}
	}
	return table.Render()
}

func formatJSON[T any](w io.Writer, results []T) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func formatCSV[T any](w io.Writer, results []T) error {
	headers, rows := extractHeadersAndRows(results)

	cw := csv.NewWriter(w)
	if err := cw.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// extractHeadersAndRows converts a slice of structs to string headers and rows.
func extractHeadersAndRows[T any](results []T) ([]string, [][]string) {
	var zero T
	t := reflect.TypeOf(zero)

	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(t.Field(i).Name)
		}
		headers = append(headers, tag)
	}

	var rows [][]string
	for _, item := range results {
		v := reflect.ValueOf(item)
		var row []string
		for i := 0; i < v.NumField(); i++ {
			row = append(row, fieldToString(v.Field(i)))
		}
		rows = append(rows, row)
	}
	return headers, rows
}

func fieldToString(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		return fmt.Sprintf("%v", v.Elem().Interface())
	}
	return fmt.Sprintf("%v", v.Interface())
}
