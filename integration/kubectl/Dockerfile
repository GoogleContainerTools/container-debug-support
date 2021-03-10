# Simple multi-platform image that includes kubectl and curl
FROM --platform=$BUILDPLATFORM curlimages/curl
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# curlimages/curl runs as curl-user and cannot install into /usr/bin
USER root
ADD https://dl.k8s.io/release/v1.20.0/bin/$TARGETOS/$TARGETARCH/kubectl /usr/bin
RUN chmod a+x /usr/bin/kubectl
