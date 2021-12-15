FROM registry.ci.openshift.org/openshift/openshift-serverless-nightly:serverless-operator-src

ENV BASE=github.com/openshift-knative/serverless-operator
WORKDIR ${GOPATH}/src/${BASE}

COPY . .

COPY test_runner.go /go/src/knative.dev/eventing/test/lib/
COPY client.go /go/src/knative.dev/eventing/test/lib/
