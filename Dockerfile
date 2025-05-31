# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Create final image
FROM alpine:latest

# Install Chrome dependencies and tools
RUN apk update && apk add --no-cache \
    chromium \
    chromium-chromedriver \
    curl \
    ca-certificates \
    tzdata \
    dbus \
    ttf-freefont \
    && mkdir -p /var/run/dbus \
    && dbus-daemon --system \
    && mkdir -p /home/appuser/.cache/rod/browser \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup \
    && chown -R appuser:appgroup /home/appuser

# Set Chrome environment variables
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/ \
    CHROMIUM_FLAGS="--headless --disable-gpu --no-sandbox --disable-dev-shm-usage --disable-software-rasterizer --disable-dbus --no-zygote --no-first-run --disable-extensions"

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Set ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set environment variables with defaults
ENV PORT=8080 \
    CHECK_INTERVAL=60

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 