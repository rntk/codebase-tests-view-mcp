# Stage 1: Build React frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copy frontend files
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.25-alpine AS go-builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code and built frontend
COPY . .
COPY --from=frontend-builder /frontend/build ./internal/api/frontend/build

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 3: Runtime image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=go-builder /server /app/server

# Create data directory for metadata
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV METADATA_PATH=/data/metadata.json

# Run the server
CMD ["/app/server", "-port", "8080", "-metadata", "/data/metadata.json"]
