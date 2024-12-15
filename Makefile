BIN_AIR = $(shell go env GOPATH)/bin/air
BIN_PKGSITE = $(shell go env GOPATH)/bin/pkgsite
BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint

GO_SOURCES=$(shell find . -name "*.go")

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
	go test ./libs/... ./tools/...

.PHONY: e2e
e2e:
	make e2e-increment-release-version

.PHONY: start-e2e-environment
start-e2e-environment:
	# admin UI 8081
	# fake server 8080
	docker compose up -d smocker
	sh wait-until-http-health.sh http://localhost:8081/version

.PHONY: stop-e2e-environment
stop-e2e-environment:
	docker compose down -v smocker

.PHONY: clean
clean:
	rm -fr dist/
	rm -rf cov/

include Makefile.tools.release.mk
