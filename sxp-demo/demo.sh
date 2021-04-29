#!/bin/bash

section() {
    DOLLAH_COLOR=green DOLLAH_DISPLAY=yes DOLLAH_PROMPT="" ./hack/dollah.py "###" "$@" "###"
}

say() {
    DOLLAH_COLOR=blue DOLLAH_DISPLAY=yes DOLLAH_PROMPT="" ./hack/dollah.py "#" "$@"
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

all_done() {
    PS1='[done]$ ' bash --norc
}

if [ -z "$IMG" ]; then
    echo "IMG is unset!"
    exit 1
fi
if [ -z "$TAG" ]; then
    echo "IMG is unset!"
    exit 1
fi

export DOLLAH_COLOR=white

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


if ! kubectl -n "$ns" exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.$ns.svc.cluster.local/Users" -c "ls" >/dev/null ; then
    err_users=true
fi

if ! kubectl -n "$ns" exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.domain1.sink.test/Users" -c "ls" >/dev/null ; then
    err_archive=true
fi

if [ "$err_users" = "true" -o "$err_archive" = "true" ]; then
    say "we can dive deeper"
else
    say "we can now use the shares - start with smbclient in cluster"
    demo kubectl -n "$ns" exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.$ns.svc.cluster.local/Users"
fi

manual

all_done
