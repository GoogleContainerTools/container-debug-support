#!/bin/sh

set -euo pipefail

echo ">> Building..."
skaffold build -p integration --skip-tests

echo ">> Launching test jobs and pods"
skaffold run -p integration --tail &
skaffoldPid=$!

trap "echo '>> Tearing down test jobs'; kill $skaffoldPid; skaffold delete -p integration" 0 1 3 15

echo ">> Waiting for test jobs to start"
sleep 5
echo ">> Monitoring for test job completion"
kubectl wait --for=condition=complete job.batch \
    -l project=container-debug-support,type=integration-test
