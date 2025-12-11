# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a static binary (disable CGO for portability)
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Stage 2: Runtime
# We use 'gcr.io/distroless/static-debian13' for maximum security (no shell, no package manager)
# or just 'alpine' if you want debugging tools.
FROM gcr.io/distroless/static-debian13

WORKDIR /

# Create non-root user (UID 1000, GID 1000)
USER 1000:1000

COPY --from=builder /app/server /server

# Expose port
EXPOSE 8000

# Run
CMD ["/server"]