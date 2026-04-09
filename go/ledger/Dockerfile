# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o ledger ./cmd/ledger

# ---- Run stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/ledger /ledger

EXPOSE 8080

ENTRYPOINT ["/ledger"]
