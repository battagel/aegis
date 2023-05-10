# First stage: Build
FROM golang:1.19 AS builder

# Set the build directory
WORKDIR /build

ADD go.mod go.sum Makefile config.env clamd.conf ./
ADD ./internal ./internal
ADD ./pkg ./pkg
ADD ./cmd ./cmd

RUN make build

# Second stage: Run
FROM debian:bullseye-slim

# Install clamdscan and clean image
RUN apt-get update && \
    apt-get install -y clamdscan && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy the Aegis binary
COPY --from=builder /build/bin/aegis /usr/local/bin/aegis

# Copy the Aegis config file
COPY --from=builder /build/config.env /app/config.env
COPY --from=builder /build/clamd.conf /app/clamd.conf

# Set the working directory
WORKDIR /app

# Expose the Aegis port
EXPOSE 8080

# Start ClamAV daemon and Aegis
CMD ["/usr/local/bin/aegis"]
