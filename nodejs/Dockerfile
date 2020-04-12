FROM golang:1.14.1 as build
COPY . .
ENV GOPATH=
RUN go build -o node wrapper.go

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM busybox
# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=build /go/node nodejs/bin/
