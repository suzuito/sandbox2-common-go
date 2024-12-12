BIN_AIR = $(shell go env GOPATH)/bin/air
BIN_PKGSITE = $(shell go env GOPATH)/bin/pkgsite
BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint

.PHONY: mac-init
mac-init:

$(BIN_GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

$(BIN_AIR):
	go install github.com/air-verse/air

$(BIN_PKGSITE):
	go install golang.org/x/pkgsite/cmd/pkgsite

.PHONY: lint
lint: $(BIN_GOLANGCI_LINT)
	$(BIN_GOLANGCI_LINT) run ./...

.PHONY: godoc
godoc: $(BIN_AIR) $(BIN_PKGSITE)
	$(BIN_AIR) -c air.godoc.toml

.PHONY: test
test:
	go test ./...

.PHONY: build-increment-release-prd
build-increment-release-prd:
	mkdir -p dist/prd && go build -o dist/prd/increment-release-version tools/release/cmd/increment-release-version/*.go

.PHONY: clean
clean:
	rm -fr dist/
