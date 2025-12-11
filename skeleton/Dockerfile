# Build stage
FROM golang:{{values.goVersion}}-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o ${{values.name}} \
    ./cmd/${{values.name}}

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 ${{values.name}} && \
    adduser -D -u 1000 -G ${{values.name}} ${{values.name}}

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/${{values.name}} /usr/local/bin/${{values.name}}

# Copy example config
COPY --chown=${{values.name}}:${{values.name}} .env.example /app/.env.example

# Switch to non-root user
USER ${{values.name}}

# Set entrypoint
ENTRYPOINT ["${{values.name}}"]
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="${{values.name}}"
LABEL org.opencontainers.image.description="${{values.description}}"
LABEL org.opencontainers.image.source="https://github.com/fast-ish/${{values.name}}"
LABEL org.opencontainers.image.vendor="fast-ish"
LABEL org.opencontainers.image.version="${VERSION}"
