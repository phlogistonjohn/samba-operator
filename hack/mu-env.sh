#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd)"
set -xe

cd "$SCRIPT_DIR/.."
ns=samba-operator-system

# test env
t=./tests/files
kubectl -n $ns apply \
    -f $t/data1.yaml \
    -f $t/joinsecret1.yaml \
    -f $t/client-test-pod.yaml \
    -f $t/smbsecurityconfig2.yaml

# multus
kubectl -n $ns apply -f hostdev1.yaml
