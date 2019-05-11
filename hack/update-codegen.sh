#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CUSTOM_HEADER=${PWD}/hack/boilerplate.go.txt
cd ${PWD}/vendor/k8s.io/code-generator

./generate-groups.sh \
  client,lister,informer \
  github.com/banzaicloud/kafka-operator/pkg/client \
  github.com/banzaicloud/kafka-operator/pkg/apis \
  banzaicloud:v1alpha1 \
  --go-header-file ${CUSTOM_HEADER}
