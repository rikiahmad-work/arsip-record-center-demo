# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/server/main.go

# Production stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta

# Create storage directory for documents
RUN mkdir -p /app/storage/dokumen

COPY --from=builder /app/bin/server /app/server

EXPOSE 8000

CMD ["/app/server"]
