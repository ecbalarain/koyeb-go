# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from build stage
COPY --from=build /app/main .
COPY --from=build /app/migrate .

# Copy migrations
COPY --from=build /app/migrations ./migrations

# Copy frontend static files (for serving admin pages and products.json)
COPY --from=build /app/cloudflare-pages-frontend ./cloudflare-pages-frontend

# Copy startup script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Expose port (Koyeb will override this with PORT env var)
EXPOSE 8080

# Use entrypoint script to run migrations then start app
ENTRYPOINT ["./docker-entrypoint.sh"]