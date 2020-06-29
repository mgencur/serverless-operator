#!/usr/bin/env bash

export PERIODICS=${PERIODICS:-false}

export IMAGE=registry.svc.ci.openshift.org/openshift/openshift-serverless-nightly:serverless-operator-src
export SERVERLESS_HOME=/go/src/github.com/openshift-knative/serverless-operator
export KNATIVE_SERVING_HOME=/go/src/knative.dev/serving
export KNATIVE_EVENTING_HOME=/go/src/knative.dev/eventing

function print_stats {
  local tests=$1
  local maxage=$(( $2 * 24 )) #convert days to hours
  for test in $tests; do
    stats=$(curl -s "https://search.apps.build01.ci.devcluster.openshift.com/?search=${test}&maxAge=${maxage}h&context=2&type=junit&name=${jobname}&maxMatches=3&maxBytes=20971520&groupBy=job" | grep "Across" | awk -F "em" '{ print $3 }' | awk -F ">" '{ print $2 }' | awk -F "<" '{ print $1 }' | awk -F "and 100.00%" '{ print $1 }')
    if [[ -n "$stats" ]]; then
      local stats_formatted=$(echo $stats | sed "s/\(.*\),\(.*\)/\2 (\1)/") #move important stuff to the beginning
      echo "$test: $stats_formatted"
    fi
  done
}

serverless_tests=$(podman run -w "$SERVERLESS_HOME" "$IMAGE" go test -tags=e2e  \
                  ./test/e2e ./test/servinge2e  \
                  --kubeconfigs /dev/null | \
                  grep "^--- FAIL:" | awk -F "FAIL: " '{ print $2}' | awk '{ print $1 }' | sort | uniq)

serving_tests=$(podman run -w "$KNATIVE_SERVING_HOME" "$IMAGE" go test -tags=e2e  \
                ./test/e2e ./test/conformance/api/... ./test/conformance/runtime/...  \
                --kubeconfig /dev/null | \
                grep "^--- FAIL:" | awk -F "FAIL: " '{ print $2}' | awk '{ print $1 }' | sort | uniq)

eventing_tests=$(podman run -w "$KNATIVE_EVENTING_HOME" "$IMAGE" go test -tags=e2e  \
                ./test/e2e \
                --kubeconfig /dev/null | \
                grep "^--- FAIL:" | awk -F "FAIL: " '{ print $2}' | awk '{ print $1 }' | sort | uniq)

for maxdays in 1 3 7; do
  printf "\n%s\n" "###### Stats for last $maxdays days ######"

  for ocp in "4.3" "4.4" "4.5"; do
    if [[ "$PERIODICS" == "true" ]]; then
      jobname="^periodic-ci-openshift-knative-serverless-operator-master-${ocp}-e2e-aws*"
    else
      jobname="^pull-ci-openshift-knative-serverless-operator-master-${ocp}-e2e-aws*"
    fi

    printf "\n%s\n" "=== OCP ${ocp} ==="

    printf "\n%s\n" "--- Serverless operator tests ---"

    print_stats "$serverless_tests" $maxdays

    printf "\n%s\n" "--- Knative Serving tests ---"

    print_stats "$serving_tests" $maxdays

    printf "\n%s\n" "--- Knative Eventing tests ---"

    print_stats "$eventing_tests" $maxdays
  done
done
