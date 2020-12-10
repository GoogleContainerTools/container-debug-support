#!/bin/sh
skaffold build -p prod --default-repo gcr.io/k8s-skaffold/debug-support-helpers
skaffold build -p prod --default-repo gcr.io/gcp-dev-tools/duct-tape
