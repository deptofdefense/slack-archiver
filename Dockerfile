# build stage
FROM golang:alpine3.13 AS builder

RUN apk update && apk add --no-cache git make gcc g++ ca-certificates && update-ca-certificates

WORKDIR /src

COPY . .

RUN rm -f bin/gox bin/slack-archiver_linux_amd64 && make bin/gox && bin/gox \
-cgo \
-os="linux" \
-arch="amd64" \
-ldflags "-s -w" \
-output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" \
github.com/deptofdefense/slack-archiver/cmd/slack-archiver

# final stage
FROM alpine:3.15

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /src/bin/slack-archiver_linux_amd64 /bin/slack-archiver

ENTRYPOINT ["/bin/slack-archiver"]
