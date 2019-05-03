FROM golang:1.11 as delve
RUN curl --location --output delve-1.2.0.tar.gz https://github.com/go-delve/delve/archive/v1.2.0.tar.gz \
 && tar xzf delve-1.2.0.tar.gz
RUN cd delve-1.2.0/cmd/dlv && go install

#
#FROM python:2.7 as python27
#RUN curl --location --output ptsvd.zip https://github.com/Microsoft/ptvsd/archive/v4.2.4.zip \
# && unzip ptsvd.zip \
# && cd ptsvd-4.2.4 \
# && 

# The debian image is about 95MB bigger
FROM gcr.io/distroless/base:debug
SHELL ["/busybox/sh","-c"]
WORKDIR /duct-tape
COPY --from=delve /go/bin/dlv .
COPY install.sh /
ENTRYPOINT ["/busybox/sh", "/install.sh"]
