FROM golang:1.21-alpine AS builder-base

WORKDIR /workdir
RUN mkdir /builddir

# Certificates

RUN apk --update upgrade
RUN apk add ca-certificates
RUN update-ca-certificates

# Build and include grpc-health-probe
RUN env GOBIN=/builddir go install github.com/grpc-ecosystem/grpc-health-probe@latest

# Copy across
COPY go.work go.work.sum ./
COPY modules/worker/go.mod modules/worker/go.sum ./modules/worker/
COPY modules/frontdoor/go.mod modules/frontdoor/go.sum ./modules/frontdoor/
COPY modules/service-pkg/go.mod modules/service-pkg/go.sum ./modules/service-pkg/
COPY modules/testing-tools/go.mod modules/testing-tools/go.sum ./modules/testing-tools/
COPY test-harness/go.mod test-harness/go.sum ./test-harness/

# Install go application dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

FROM builder-base AS builder-instrumented

ARG entrypoint

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
  CGO_ENABLED=0 go build -cover -v -o /builddir/service ./$entrypoint

FROM builder-base as builder

ARG entrypoint

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -v -o /builddir/service ./$entrypoint

FROM gcr.io/distroless/static as runtime-base

COPY --from=builder-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder-base /builddir/grpc-health-probe /bin

# non root
USER 65532

FROM runtime-base as runtime-instrumented

COPY --from=builder-instrumented /builddir/service /bin

ENTRYPOINT ["/bin/service"]

FROM runtime-base AS runtime

COPY --from=builder /builddir/service /bin

ENTRYPOINT ["/bin/service"]
