#!/bin/bash
# Run Integration Tests
# Can be called either from the top-level directory or from any of
# the language-specific images.
#
# $ sh run-its.sh
# $ (cd go; sh ../run-its.sh)
#
# Integration tests are set up as a set of Jobs.  This script launches
# a set of Pods and the Jobs, and then waits for the Jobs to complete.

set -euo pipefail

countTestJobs() {
  kubectl get job.batch -o name -l project=container-debug-support,type=integration-test \
  | wc -l
}

echo ">> Building test images [$(date)]"
skaffold build -p integration

echo ">> Launching test jobs and pods [$(date)]"
skaffold run -p integration --tail &
skaffoldPid=$!

trap "echo '>> Tearing down test jobs [$(date)]'; kill $skaffoldPid; skaffold delete -p integration" 0 1 3 15

echo ">> Waiting for test jobs to start [$(date)]"
jobcount=0
while [ $jobcount -eq 0 -o $jobcount -ne $(countTestJobs) ]; do
    jobcount=$(countTestJobs)
    sleep 5
done

echo ">> Monitoring for test job completion [$(date)]"
kubectl wait --for=condition=complete job.batch \
    -l project=container-debug-support,type=integration-test
