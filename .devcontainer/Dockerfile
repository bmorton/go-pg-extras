FROM mcr.microsoft.com/devcontainers/base:ubuntu

RUN apt-get update && apt-get install -y --no-install-recommends postgresql-client \
    && rm -rf /var/lib/apt/lists/*

ARG TARGETARCH
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-${TARGETARCH} \
    -o /usr/local/bin/tailwindcss && chmod +x /usr/local/bin/tailwindcss
