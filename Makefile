BIN_AIR = $(shell go env GOPATH)/bin/air
BIN_GODOC = $(shell go env GOPATH)/bin/godoc
BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint

.PHONY: mac-init
mac-init:

$(BIN_GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

$(BIN_AIR):
	go install github.com/air-verse/air@v1.61.1

$(BIN_GODOC):
	go install golang.org/x/tools/cmd/godoc@v0.28.0

.PHONY: lint
lint:
	$(BIN_GOLANGCI_LINT) run ./...

.PHONY: godoc
godoc: $(BIN_AIR) $(BIN_GODOC)
	$(BIN_AIR) -c air.godoc.toml

.PHONY: test
test:
	go test ./...
