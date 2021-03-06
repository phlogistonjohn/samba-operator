CONTAINER := quay.io/obnox/samba-operator:v0.0.1

OUTPUT ?= build/_output
TOOLS_DIR ?= build/_tools

export OPERATOR_SDK_VERSION ?= v0.17.1
export OPERATOR_SDK ?= $(TOOLS_DIR)/operator-sdk-$(OPERATOR_SDK_VERSION)

export GO111MODULE=on
# GOROOT is needed for the operator-sdk to work
export GOROOT ?= $(shell go env GOROOT)
# use GOPROXY by default to speed up dependency operations
export GOPROXY ?= https://proxy.golang.org

all: build

operator-sdk: $(OPERATOR_SDK)
.PHONY: operator-sdk

$(OPERATOR_SDK):
	@echo "Ensuring operator-sdk"
	hack/ensure-operator-sdk.sh
	@echo "[OK] operator-sdk"

generate: generate.crds generate.k8s
.PHONY: generate

generate.k8s: $(OPERATOR_SDK)
	$(OPERATOR_SDK) generate k8s
	@echo "[OK] generate.k8s"
.PHONY: generate.k8s

generate.crds: $(OPERATOR_SDK)
	$(OPERATOR_SDK) generate crds
	@echo "[OK] generate.crds"
.PHONY: generate.crds

build: $(OPERATOR_SDK) generate
	$(OPERATOR_SDK) build $(CONTAINER)
	@echo "[OK] build"
.PHONY: build

clean:
	rm -rf $(OUTPUT)
	rm -f go.sum
	@echo "[OK] clean"
.PHONY: clean

realclean: clean
	rm -rf $(TOOLS_DIR)
	@echo "[OK] realclean"
.PHONY: realclean

push: build
	docker push $(CONTAINER)
	@echo "[OK] push"
.PHONY: push

install:
	./deploy/install.sh
.PHONY: install

uninstall:
	./deploy/uninstall.sh
.PHONY: uninstall
