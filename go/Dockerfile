FROM golang:1.14.1 as delve

RUN curl --location --output delve-1.4.0.tar.gz https://github.com/go-delve/delve/archive/v1.4.0.tar.gz \
  && tar xzf delve-1.4.0.tar.gz \
  && mv delve-1.4.0 delve-source

# Patch delve to change default for --only-same-user to false
# Required as `kubectl port-forward` to dlv port is refused.
# We must install patch(1) to apply the patch.
COPY delve-only-same-user.patch .
RUN apt-get update && apt-get install -y --no-install-recommends \
  patch \
  && patch -p0 -d delve-source < delve-only-same-user.patch

# Produce an as-static-as-possible dlv binary to work on musl and glibc
RUN cd delve-source && CGO_ENABLED=0 go build -o /go/dlv -ldflags '-s -w -extldflags "-static"' ./cmd/dlv/

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM busybox
# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=delve /go/dlv go/bin/
