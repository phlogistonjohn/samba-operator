#!/bin/bash

section() {
    echo ""
    DOLLAH_COLOR=green DOLLAH_DISPLAY=yes DOLLAH_PROMPT="" ./hack/dollah.py "###" "$@" "###"
}

say() {
    echo ""
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
    WHITE="\033[1;37;40m"
    NORMAL="\033[0;37;40m"
    PS1="${WHITE}[show]\$${NORMAL} " bash --norc
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

must demo kubectl get namespaces

# deploy the operator (manually, via makefile)
must demo make deploy

must demo kubectl get namespaces

must demo kubectl config set-context --current --namespace=samba-operator-system

must demo kubectl get pods

# deploy the smb client pod, etc.
must demo kubectl apply -f ./sxp-demo/client.yaml

demo kubectl get pods -w


section "Deploying the Custom Resources - Without Active Directory"

# show some files
manual

must demo kubectl apply -f ./sxp-demo/ugsec1.yaml
must demo kubectl apply -f ./sxp-demo/ugshare1.yaml

demo kubectl get pods -w

say "we can now access the share"
demo kubectl exec -it smbclient -- smbclient -U 'sambauser%letm31n_pls' "//first.$ns.svc.cluster.local/First Share"

say "what do we see within the container?"
demo kubectl exec -it deploy/first -- bash

section "Deploying the Custom Resources - With Active Directory"

# show some files
manual


must demo kubectl apply -f ./sxp-demo/adsecurity.yaml
must demo kubectl apply -f ./sxp-demo/adshare1.yaml
must demo kubectl apply -f ./sxp-demo/adshare2.yaml

demo kubectl get pods -w


if ! kubectl exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.$ns.svc.cluster.local/Users" -c "ls" >/dev/null ; then
    err_users=true
fi

if ! kubectl exec -it smbclient -- smbclient -U 'DOMAIN1/bwayne%1115Rose.' "//users.domain1.sink.test/Users" -c "ls" >/dev/null ; then
    err_archive=true
fi

if [ "$err_users" = "true" -o "$err_archive" = "true" ]; then
    say "we can dive deeper"
else
    say "now to copy some files - using smbclient"
    demo kubectl exec -it smbclient -- smbclient -U 'DOMAIN1/ckent%1115Rose.' "//users.$ns.svc.cluster.local/Users"
fi

say "time to access the shares from outside the cluster"
manual

all_done
