FROM golang:1.14.1 as delve
RUN curl --location --output delve-1.3.1.tar.gz https://github.com/go-delve/delve/archive/v1.3.1.tar.gz \
  && tar xzf delve-1.3.1.tar.gz
# Produce an as-static-as-possible dlv binary to work on musl and glibc
RUN cd delve-1.3.1 && CGO_ENABLED=0 go build -o /go/dlv -ldflags '-extldflags "-static"' ./cmd/dlv/

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM busybox
# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=delve /go/dlv go/bin/
