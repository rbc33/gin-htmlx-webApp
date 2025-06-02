
FROM ubuntu:latest


ENV TZ=Europe/London
ENV DEBIAN_FRONTEND="noninteractive"

ARG TARGETARCH=amd64
ENV PATH="${PATH}:/usr/local/go/bin"

#
LABEL maintainer="ricardobenthem@gmail.com"

RUN apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests \
    ca-certificates \
    git \
    make \
    pkg-config \
    wget \
    && wget "https://go.dev/dl/go1.22.0.linux-${TARGETARCH}.tar.gz" -P golang-install \
    && cd golang-install \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf "go1.22.0.linux-${TARGETARCH}.tar.gz" \
    && GOBIN=/usr/local/bin go install github.com/pressly/goose/v3/cmd/goose@v3.18.0 \
    && GOBIN=/usr/local/bin go install github.com/a-h/templ/cmd/templ@v0.2.543 \
    && apt-get remove -y wget \
    && apt-get clean \
    && apt-get autoremove -y \
    && apt-get autoremove --purge -y \
    && go clean --cache \
    && go clean --modcache \
    && go clean --testcache \
    && go clean --fuzzcache \
    && rm -rf golang-install \
    && rm -rf /var/lib/apt/lists/* \