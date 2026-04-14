# --- Stage 1: Build ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# --- Stage 2: Final ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary and migrations from builder
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 3000

CMD ["./api"]
