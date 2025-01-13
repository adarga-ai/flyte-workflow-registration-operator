# Build the manager binary
FROM golang:1.23 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG FLYTECTL_VERSION

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/

# Install curl and necessary utilities for checksum validation
RUN apt-get update && apt-get install -y curl gnupg jq

# Download and validate the installation script
RUN curl -sSL "https://ctl.flyte.org/install" -o install.sh \
    && echo 'd98173e3276c2ff37daf37c130137a3df00311a9c313c40d9f9da66dfe43f17a install.sh' | sha256sum -c - \
    && chmod +x install.sh

# Execute the installation script with the desired version
RUN ./install.sh ${FLYTECTL_VERSION}

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:debug-nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
# Ensure flytectl is available in the distroless image
COPY --from=builder /workspace/bin/flytectl /usr/local/bin/flytectl
USER 65532:65532

ENTRYPOINT ["/manager"]
