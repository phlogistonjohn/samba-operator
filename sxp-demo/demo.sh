#!/bin/bash

section() {
    echo "###" "$@" "###"
}

say() {
    demo "echo" "$@"
}

demo() {
    ./hack/dollah.py "$@"
    "$@"
}

must() {
    "$@" || {
        echo "ERROR"
        exit 1
    }
}

manual() {
    PS1='[show]$ ' bash --norc
}

ns=samba-operator-system

section "Deploying the operator"

# deploy the operator (manually, via makefile)
must demo make deploy

must demo kubectl -n "$ns" get pods

# deploy the smb client pod, etc.
must demo kubectl -n "$ns" apply -f ./sxp-demo/client.yaml

must demo kubectl -n "$ns" get pods


section "Deploying the Custom Resources - Without Active Directory"

# show some files
manual

must demo kubectl -n "$ns" apply -f ./sxp-demo/ugsec1.yaml
must demo kubectl -n "$ns" apply -f ./sxp-demo/ugshare1.yaml

demo kubectl -n "$ns" get pods -w

say "let us access the share"
demo kubectl -n "$ns" exec -it smbclient -- smbclient -U 'sambauser%letm31n_pls' "//first.$ns.svc.cluster.local/First Share"

say "what do we see within the container?"
demo kubectl -n "$ns" exec -it deploy/first -- bash

section "Deploying the Custom Resources - With Active Directory"

# show some files
manual


must demo kubectl -n "$ns" apply -f ./sxp-demo/adsecurity.yaml
must demo kubectl -n "$ns" apply -f ./sxp-demo/adshare1.yaml
must demo kubectl -n "$ns" apply -f ./sxp-demo/adshare2.yaml

demo kubectl -n "$ns" get pods -w


must demo kubectl -n "$ns" exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.$ns.svc.cluster.local/Users" -c "ls"
must demo kubectl -n "$ns" exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.domain1.sink.test/Users" -c "ls"

manual
