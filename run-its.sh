#!/bin/bash
# run integration tests
# any arguments are passed onto skaffold
#
# $ sh run-its.sh -p local
#

set -euo pipefail

echo ">> Building test images [$(date)]"
skaffold build -p integration "$@"

echo ">> Launching test jobs and pods [$(date)]"
skaffold run -p integration --tail "$@" &
skaffoldPid=$!

trap "echo '>> Tearing down test jobs [$(date)]'; kill $skaffoldPid; skaffold delete -p integration" 0 1 3 15

echo ">> Waiting for test jobs to start [$(date)]"
while [ $(kubectl get job.batch -o name -l project=container-debug-support,type=integration-test | wc -l) -lt 1 ]; do
    sleep 5
done
# Do an additional wait to ensure all tests are up and running as
# `kubectl wait` doesn't refresh and will miss late appearing resources
sleep 5
kubectl get job.batch -o name -l project=container-debug-support,type=integration-test

echo ">> Monitoring for test job completion [$(date)]"
kubectl wait --for=condition=complete job.batch \
    -l project=container-debug-support,type=integration-test
