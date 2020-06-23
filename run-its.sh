#!/bin/bash

set -euo pipefail

echo ">> Building test images [$(date)]"
skaffold build -p integration --skip-tests

echo ">> Launching test jobs and pods [$(date)]"
skaffold run -p integration --tail &
skaffoldPid=$!

trap "echo '>> Tearing down test jobs [$(date)]'; kill $skaffoldPid; skaffold delete -p integration" 0 1 3 15

echo ">> Waiting for test jobs to start [$(date)]"
while [ $(kubectl get job.batch -o name | wc -l) -eq 0 ]; do
    sleep 5
done
echo ">> Monitoring for test job completion [$(date)]"
kubectl wait --for=condition=complete job.batch \
    -l project=container-debug-support,type=integration-test
