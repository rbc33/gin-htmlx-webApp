FROM ubuntu:latest

ENV TZ=Europe/London
ENV DEBIAN_FRONTEND="noninteractive"

ARG TARGETARCH=amd64
ARG GO_VERSION="1.23.2"      
ARG GOOSE_VERSION="3.18.0"
ARG TEMPL_VERSION="0.2.543"
ARG GOLANGCI_LINT_VERSION="v1.58.0" 
ENV PATH="${PATH}:/usr/local/go/bin:/usr/local/bin" 

#
LABEL maintainer="ricardobenthem@gmail.com"

RUN apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests \
    ca-certificates \
    git \
    make \
    pkg-config \
    wget \
    && wget "https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" -P /tmp/ \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" \
    && rm -f "/tmp/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" \
    \
    && GOBIN=/usr/local/bin go install github.com/pressly/goose/v3/cmd/goose@v${GOOSE_VERSION} \
    && GOBIN=/usr/local/bin go install github.com/a-h/templ/cmd/templ@v${TEMPL_VERSION} \
    && GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} \
    \
    && apt-get clean \
    && apt-get autoremove -y \
    && apt-get autoremove --purge -y \
    && go clean --cache \
    && go clean --modcache \
    && go clean --testcache \
    && go clean --fuzzcache 


RUN if [ "$TARGETARCH" = "amd64" ]; then \
    wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-linux-x64 -O /usr/local/bin/tailwindcss && chmod +x /usr/local/bin/tailwindcss; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
    wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-linux-arm64 -O /usr/local/bin/tailwindcss && chmod +x /usr/local/bin/tailwindcss; \
    else \
    echo "Unsupported architecture: $TARGETARCH"; \
    exit 1; \
    fi
RUN apt-get remove -y wget \
    && rm -rf /var/lib/apt/lists/*  
