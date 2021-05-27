ARG GOVERSION=1.16
FROM --platform=$BUILDPLATFORM golang:${GOVERSION} as build
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY . .
# Produce an as-static-as-possible dlv binary to work on musl and glibc
RUN GOPATH="" CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o node -ldflags '-s -w -extldflags "-static"' wrapper.go

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM busybox
# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=build /go/node nodejs/bin/
