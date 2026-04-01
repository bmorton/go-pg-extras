FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY pgextras /usr/local/bin/pgextras

ENTRYPOINT ["pgextras"]
