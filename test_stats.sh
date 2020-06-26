#!/usr/bin/env bash

# This script must be run in the Serving or Eventing home directory to have access to tests.

export PERIODICS=${PERIODICS:-false}

# Serving
tests=$(go test -tags=e2e  \
        ./test/e2e ./test/conformance/api/... ./test/conformance/runtime/...  \
        --kubeconfig /dev/null | \
        grep "^--- FAIL:" | awk -F "FAIL: " '{ print $2}' | awk '{ print $1 }')

# Eventing
# tests=$(go test -tags=e2e  \
#         ./test/e2e  \
#         --kubeconfig /dev/null | \
#         grep "^--- FAIL:" | awk -F "FAIL: " '{ print $2}' | awk '{ print $1 }')

for ocp in "4.3" "4.4" "4.5"; do

  if [[ "$PERIODICS" == "true" ]]; then
    jobname="^periodic-ci-openshift-knative-serverless-operator-master-${ocp}-e2e-aws*"
  else
    jobname="^pull-ci-openshift-knative-serverless-operator-master-${ocp}-e2e-aws*"
    #jobname="^pull-ci-openshift-knative-serverless-operator-master-4.3-upgrade-tests*"
  fi

  echo "=== OCP ${ocp} ==="

  for test in $tests; do
    stats=$(curl -s "https://search.apps.build01.ci.devcluster.openshift.com/?search=${test}&maxAge=168h&context=2&type=junit&name=${jobname}&maxMatches=1&maxBytes=20971520&groupBy=job" | grep "Across" | awk -F "em" '{ print $3 }' | awk -F ">" '{ print $2 }' | awk -F "<" '{ print $1 }' | awk -F "and 100.00%" '{ print $1 }')
    if [[ -n "$stats" ]]; then
      echo "$test: $stats"
    fi
  done
done
