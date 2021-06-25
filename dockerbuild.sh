#!/bin/sh
#
# Skaffold Custom Builder that uses `docker buildx` to perform a
# multi-platform build when necessary.

export DOCKER_BUILDKIT=1

if [ "$PUSH_IMAGE" = true ]; then
    if [ -z "$PLATFORMS" ]; then
        PLATFORMS=linux/amd64,linux/arm64
    fi
    loadOrPush="--platform $PLATFORMS --push"
else
    # cannot load multiarch images into the daemon
    set -x
    exec docker buildx build "$@" --load --tag $IMAGE "$BUILD_CONTEXT"
fi

if ! docker buildx inspect skaffold-builder >/dev/null 2>&1; then
  echo ">> creating 'docker buildx' builder 'skaffold-builder'"
  # Docker 3.3.0 require creating a builder within a context
  (set -x; \
    docker context create skaffold; \
    docker buildx create --name skaffold-builder ${PLATFORMS:+--platform $PLATFORMS} skaffold)
fi

set -x
docker buildx build \
  --progress=plain \
  --builder skaffold-builder \
  $loadOrPush \
  "$@" \
  --tag $IMAGE \
  "$BUILD_CONTEXT"

