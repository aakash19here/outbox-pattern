# syntax=docker/dockerfile:1

# ---- build stage -----------------------------------------------------------
# Debian-based (not alpine) because confluent-kafka-go links against librdkafka
# via cgo, which is simplest with the distro packages.
FROM golang:1.26 AS build

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential pkg-config librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Cache module downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both service binaries from the single module
RUN CGO_ENABLED=1 go build -tags dynamic -o /out/http-server ./cmd/http-server/. \
 && CGO_ENABLED=1 go build -tags dynamic -o /out/msg-server ./cmd/msg-server/.

# ---- runtime stage ---------------------------------------------------------
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/http-server /usr/local/bin/http-server
COPY --from=build /out/msg-server /usr/local/bin/msg-server

# Which binary runs is selected per-service in compose (`command`)
CMD ["/usr/local/bin/msg-server"]