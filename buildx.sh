#!/bin/sh
#
# Skaffold Custom Builder that uses `docker buildx` to perform a
# multi-platform build.

export DOCKER_BUILDKIT=1

# The local Docker daemon which cannot load images for multiple architectures,
# so just build using normal Docker.
if [ "$PUSH_IMAGE" != true ]; then
    exec docker build --tag $IMAGE "$BUILD_CONTEXT"
fi

PLATFORMS=linux/amd64,linux/arm64

if ! docker buildx inspect skaffold-builder >/dev/null 2>&1; then
  echo ">> creating 'docker buildx' builder 'skaffold-builder'"
  # Docker 3.3.0 require creating a builder within a context
  (set -x; \
    docker context create skaffold; \
    docker buildx create --name skaffold-builder --platform $PLATFORMS skaffold)
fi

set -x
docker buildx build \
  --progress=plain \
  --builder skaffold-builder \
  --platform "$PLATFORMS" \
  --push \
  --tag $IMAGE \
  "$BUILD_CONTEXT"

