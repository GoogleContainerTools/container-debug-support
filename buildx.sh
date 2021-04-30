#!/bin/sh
#
# Skaffold Custom Builder that uses `docker buildx` to perform a
# multi-platform build.

PLATFORMS=linux/amd64,linux/arm64

if ! docker buildx inspect skaffold-builder >/dev/null 2>&1; then
  echo ">> creating "docker buildx" builder 'skaffold-builder'"
  # Docker 3.3.0 require creating a builder within a context
  docker context create skaffold
  docker buildx create --name skaffold-builder --platform $PLATFORMS skaffold
fi

loadOrPush=$(if [ "$PUSH_IMAGE" = true ]; then echo --platform $PLATFORMS --push; else echo --load; fi)

set -x
docker buildx build \
  --progress=plain \
  --builder skaffold-builder \
  $loadOrPush \
  --tag $IMAGE \
  "$BUILD_CONTEXT"

