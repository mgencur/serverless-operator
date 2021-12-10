#!/usr/bin/env bash

# For SC2164
set -e

function upstream_knative_eventing_e2e {
  logger.info 'Running eventing tests'

  export TEST_IMAGE_TEMPLATE="registry.ci.openshift.org/openshift/knative-${KNATIVE_EVENTING_VERSION}:knative-eventing-test-{{.Name}}"

  cd "${KNATIVE_EVENTING_HOME}"

  # shellcheck disable=SC1091
  source "${KNATIVE_EVENTING_HOME}/openshift/e2e-common.sh"

  logger.info 'Installing Tracing'
  install_tracing

  # Eventing E2E require the KUBECONFIG.
  # TODO: Remove this when upstream tests can use in-cluster config.
  if [[ -z "$KUBECONFIG" ]]; then
    create_cluster_admin
    KUBECONFIG="$(pwd)/kubeadmin.kubeconfig"
  fi

  # run_e2e_tests defined in knative-eventing
  logger.info 'Starting eventing e2e tests'
  run_e2e_tests

  # run_conformance_tests defined in knative-eventing
  logger.info 'Starting eventing conformance tests'
  run_conformance_tests
}
