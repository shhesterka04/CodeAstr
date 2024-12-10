FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/glaphyra ./cmd/glaphyra

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/glaphyra .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

CMD ["./glaphyra"]