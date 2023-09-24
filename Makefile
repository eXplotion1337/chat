SERVICE_NAME ?= chat
REFLEX ?= github.com/cespare/reflex
SEP			 ?= "========================================================"

.PHONY: dev
dev:
	go get github.com/cespare/reflex
	go run $(REFLEX) -R "\\.idea|vendor" -r "\\.go" -s -- sh -c "GO111MODULE=on go run --race ./cmd/..."

################################################################################################################.PHONY: dev
