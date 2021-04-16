FROM golang:1.16.3-alpine3.13 as builder

COPY go.mod /src/url-collector/
WORKDIR /src/url-collector/
RUN go mod download
COPY . /src/url-collector/
RUN go build -o build/url-collector

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /src/url-collector/build/url-collector /usr/bin/url-collector
EXPOSE ${PORT} ${PORT}
ENTRYPOINT ["/usr/bin/url-collector"]