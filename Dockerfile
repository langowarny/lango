FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with CGO disabled (pure Go SQLite driver)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o lango ./cmd/lango

# Runtime image
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    chromium \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set Chrome path for rod
ENV ROD_BROWSER=/usr/bin/chromium

# Create user and group
RUN groupadd -r lango && useradd -r -g lango -m -d /home/lango lango

WORKDIR /app

COPY --from=builder /app/lango /usr/local/bin/lango

# Create data directory and set permissions
RUN mkdir -p /data && chown -R lango:lango /data && chmod 700 /data

# Switch to non-root user
USER lango

EXPOSE 18789

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:18789/health || exit 1

ENTRYPOINT ["lango"]
CMD ["serve"]
