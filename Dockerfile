ARG GO_VERSION=1.25
ARG OS_RELEASE=bookworm

FROM golang:${GO_VERSION}-${OS_RELEASE} AS base

ARG GO_LINT_VERSION=v1.51.1
ARG APP_TYPE=http

# Create non root user
RUN useradd -u 10001 nonroot

# Environment variables
ENV GONOSUMDB=github.com/OmisePayments

RUN apt-get update && \
    apt-get -y --no-install-recommends install \
    ca-certificates \
    tzdata \
    git \
    && apt-get clean

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

#
# Testing build
#
FROM base AS build-test

COPY .buildkite/tools/test.txt .
RUN cat test.txt | xargs -tI % go install %

FROM build-test AS test

#
# Tools build
#
FROM base AS build-tools

ARG PROTO_PATH=/usr/local/include/google/api/

RUN apt-get install -y protobuf-compiler && \
    go install -mod=mod \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
    mkdir -p ${PROTO_PATH} && cd ${PROTO_PATH} && \
    curl -sSf -O "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/{annotations,field_behaviour,http,httpbody}.proto" \
    && apt-get clean

ENTRYPOINT ["protoc"]

#
# Tools
#
FROM build-tools AS tools

#
# Development build
#
FROM base AS build-dev

RUN go install github.com/cespare/reflex@latest

ENTRYPOINT ["reflex", "-d", "none", "-r", "\\.go$", "-s", "--"]

#
# Development
#
FROM build-dev AS dev

#
# Production build
#
FROM base AS build

COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags '-s' -o /bin/app ./cmd/${APP_TYPE}

#
# Debug
#
FROM base AS debug

RUN go install github.com/go-delve/delve/cmd/dlv

COPY . .

#
# Release
#
FROM scratch AS release

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /bin/app /app/cmd

USER 10001:nonroot

ENTRYPOINT ["/app/cmd"]
