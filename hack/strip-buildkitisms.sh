#!/bin/sh
# Strip out any Buildkit directives from Dockerfiles.
set -e

goarch() {
  case $(uname -m) in
  x86_64) echo amd64;;
  aarch64) echo arm64;;
  *) uname -m;;
  esac
}

find . -name Dockerfile \
| while read dockerfile; do
    sed \
      -e 's@FROM --platform=[^ ]* @FROM @' \
      -e 's@ARG TARGETPLATFORM$@&=linux/'$(goarch)@ \
      -e 's@ARG TARGETOS$@&=linux@' \
      -e 's@ARG TARGETARCH$@&='$(goarch)@ \
      -e 's@ARG BUILDPLATFORM$@&=linux/'$(goarch)@ \
      -e 's@ARG BUILDOS$@&=linux@' \
      -e 's@ARG BUILDARCH$@&='$(goarch)@ \
      $dockerfile > $dockerfile.$$ \
    && mv $dockerfile.$$ $dockerfile 
done

