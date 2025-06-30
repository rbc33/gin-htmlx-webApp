FROM alpine:latest AS builder

ENV TZ=Europe/London
ARG TARGETARCH
ARG GO_VERSION="1.24.4"      
ARG GOOSE_VERSION="3.24.2"
ARG TEMPL_VERSION="0.3.865"
ARG GOLANGCI_LINT_VERSION="v1.64.8" 

# Install build dependencies
RUN apk update && apk add --no-cache \
    ca-certificates \
    git \
    make \
    pkgconf \
    wget \
    build-base \
    binutils

# Install Go
RUN if [ "$TARGETARCH" = "arm64" ]; then \
    GOARCH="arm64"; \
    elif [ "$TARGETARCH" = "amd64" ]; then \
    GOARCH="amd64"; \
    else \
    echo "Unsupported architecture: $TARGETARCH"; exit 1; \
    fi \
    && wget "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" -P /tmp/ \
    && tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-${GOARCH}.tar.gz" \
    && rm -f "/tmp/go${GO_VERSION}.linux-${GOARCH}.tar.gz"

ENV PATH="${PATH}:/usr/local/go/bin"

# Install Go tools with cleanup in same layer
RUN export CGO_ENABLED=0 \
    && go install github.com/pressly/goose/v3/cmd/goose@v${GOOSE_VERSION} \
    && go install github.com/a-h/templ/cmd/templ@v${TEMPL_VERSION} \
    && go install github.com/cosmtrek/air@v1.49.0 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} \
    && mv /root/go/bin/* /usr/local/bin/ \
    && rm -rf /root/go /root/.cache

# Download tailwindcss
RUN if [ "$TARGETARCH" = "amd64" ]; then \
    wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-linux-x64 -O /usr/local/bin/tailwindcss; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
    wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-linux-arm64 -O /usr/local/bin/tailwindcss; \
    fi \
    && chmod +x /usr/local/bin/tailwindcss

# Final stage - much smaller
FROM alpine:latest

ENV TZ=Europe/London
ENV PATH="${PATH}:/usr/local/go/bin:/usr/local/bin"

# Install only runtime dependencies
RUN apk update && apk add --no-cache \
    ca-certificates \
    make \
    git \
    bash \
    tzdata \
    && rm -rf /var/cache/apk/*

# Copy Go from builder
COPY --from=builder /usr/local/go /usr/local/go

# Copy installed tools from builder
COPY --from=builder /usr/local/bin/goose /usr/local/bin/
COPY --from=builder /usr/local/bin/templ /usr/local/bin/
COPY --from=builder /usr/local/bin/air /usr/local/bin/
COPY --from=builder /usr/local/bin/golangci-lint /usr/local/bin/
COPY --from=builder /usr/local/bin/tailwindcss /usr/local/bin/

WORKDIR /gocms


LABEL maintainer="ricardobenthem@gmail.com"