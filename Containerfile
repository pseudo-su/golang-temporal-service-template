FROM golang:1.21-alpine AS builder-base

WORKDIR /workdir
RUN mkdir /builddir

# Certificates

RUN apk --update upgrade
RUN apk add ca-certificates
RUN update-ca-certificates

# Build and include grpc-health-probe
RUN env CGO_ENABLED=0 GOBIN=/builddir go install github.com/grpc-ecosystem/grpc-health-probe@645566f

# Copy across
COPY go.mod go.sum ./

# Install go application dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

FROM builder-base AS builder

ARG entrypoint
ARG gobuildopts=

ADD . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build ${gobuildopts} -v -o /builddir/service ./$entrypoint

FROM gcr.io/distroless/static AS runtime-base

USER nonroot

COPY --from=builder-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder-base /builddir/grpc-health-probe /bin

FROM runtime-base AS runtime

USER nonroot

COPY --from=builder /builddir/service /bin

ENTRYPOINT ["/bin/service"]
