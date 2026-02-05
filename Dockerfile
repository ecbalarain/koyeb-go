# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from build stage
COPY --from=build /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]