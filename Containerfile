FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main ./cmd/shop/shop.go
RUN go test ./...

FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .
# COPY --chown=appuser:appgroup ./configs/ ./configs/

RUN chmod +x /app/main && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/main"]
