#!/bin/sh
set -eu

# note: we build the first build with --cache-artifacts=false to
# avoid cache poisoning from single-platform builds
# https://github.com/GoogleContainerTools/skaffold/issues/5504

# publish with longer image names
skaffold build -p release --default-repo gcr.io/k8s-skaffold/skaffold-debug-support --cache-artifacts=false
skaffold build -p release --default-repo gcr.io/gcp-dev-tools/duct-tape

# the github project packages is a backup location; will need to
# migrate to ghcr.io/googlecontainertools at some point
skaffold build -p release --default-repo docker.pkg.github.com/googlecontainertools/skaffold

# publish with shorter (deprecated) image names
skaffold build -p release,deprecated-names --default-repo gcr.io/k8s-skaffold/skaffold-debug-support
skaffold build -p release,deprecated-names --default-repo gcr.io/gcp-dev-tools/duct-tape
