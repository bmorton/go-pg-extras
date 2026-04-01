FROM alpine:3.21

ARG TARGETARCH

RUN apk add --no-cache ca-certificates tzdata

COPY linux/${TARGETARCH}/pgextras /usr/local/bin/pgextras

ENTRYPOINT ["pgextras"]
