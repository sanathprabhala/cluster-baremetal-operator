#!/bin/bash

# Assuming a running bare metal cluster, ensure the MAO deployment
# is scaled down before running the CBO with the required env vars set

set -x
set -o errexit
set -o nounset

source ./hack/runutils.sh

scale_down_mao

set_operator_env

export WATCH_NAMESPACE="openshift-machine-api"
export POD_NAME="cluster-baremetal-operator"
export OPERATOR_NAME="cluster-baremetal-operator"

./build/_output/bin/cluster-baremetal-operator
