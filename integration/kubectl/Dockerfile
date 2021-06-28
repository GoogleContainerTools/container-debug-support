# Simple multi-platform image that includes kubectl and curl
FROM curlimages/curl
ARG BUILDPLATFORM=linux/arm64
ARG TARGETOS=linux
ARG TARGETARCH=arm64

# curlimages/curl runs as curl-user and cannot install into /usr/bin
USER root
ADD https://dl.k8s.io/release/v1.20.0/bin/$TARGETOS/$TARGETARCH/kubectl /usr/bin
RUN chmod a+x /usr/bin/kubectl
