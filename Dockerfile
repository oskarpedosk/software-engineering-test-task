FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-w -s' -o cruder ./cmd/main.go

FROM alpine:3.20

RUN adduser -D -u 1000 appuser

WORKDIR /home/appuser
COPY --from=builder --chown=appuser:appuser /build/cruder .

USER appuser
EXPOSE 8080

CMD ["./cruder"]
