# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -tags timetzdata -o ledger ./cmd/ledger

# ---- Run stage ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/ledger /ledger

EXPOSE 9090

ENTRYPOINT ["/ledger"]
