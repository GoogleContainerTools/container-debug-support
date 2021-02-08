#!/bin/sh
set -eu

# publish with longer image names
skaffold build -p prod --default-repo gcr.io/k8s-skaffold/skaffold-debug-support
skaffold build -p prod --default-repo gcr.io/gcp-dev-tools/duct-tape

# the github project packages is a backup location; will need to
# migrate to ghcr.io/googlecontainertools at some point
skaffold build -p prod --default-repo docker.pkg.github.com/googlecontainertools/skaffold

# publish with shorter (deprecated) image names
skaffold build -p prod,deprecated-names --default-repo gcr.io/k8s-skaffold/skaffold-debug-support
skaffold build -p prod,deprecated-names --default-repo gcr.io/gcp-dev-tools/duct-tape
