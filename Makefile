BIN_AIR = $(shell go env GOPATH)/bin/air
BIN_PKGSITE = $(shell go env GOPATH)/bin/pkgsite
BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint

GO_SOURCES=$(shell find . -name "*.go")

.PHONY: mac-init
mac-init:

$(BIN_GOLANGCI_LINT): Makefile
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

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
	mkdir -p cov/ && rm -f cov/*
	go test -cover ./libs/... ./tools/... -args -test.gocoverdir=$(abspath cov)
	sh report-gocovdir.sh cov

.PHONY: e2e
e2e:
	make e2e-release-increment-release-version
	make e2e-terraform-check-terraform-rules
	make e2e-terraform-terraform_on_github_action
	make e2e-fakecmd-fakecmd

.PHONY: test
merge-test-report:
	go tool covdata percent -i=\
	cov\
	,tools/release/cov/e2e/increment-release-version\
	,tools/terraform/cov/e2e/check-terraform-rules\
	,tools/terraform/cov/e2e/terraform_on_github_action\
	,tools/fakecmd/cov/e2e/fakecmd\
	 -o=textfmt.0.txt
	grep -v 'libs/e2ehelpers/' textfmt.0.txt > textfmt.1.txt
	go tool cover -html=textfmt.1.txt -o=gocov.html
	go tool cover -func=textfmt.1.txt -o=gocovfunc.txt

.PHONY: start-e2e-environment
start-e2e-environment:
	# admin UI 8081
	# fake server 8080
	docker compose up -d smocker
	sh wait-until-http-health.sh http://localhost:8081/version

.PHONY: stop-e2e-environment
stop-e2e-environment:
	docker compose down -v smocker

.PHONY: test-local
test-local:
	make test e2e merge-test-report && sh fail-if-coverage-unsatisfied.sh 80

.PHONY: test-ci
test-ci:
	make test e2e merge-test-report

include Makefile.tools.release.mk
include Makefile.tools.terraform.mk
include Makefile.tools.fakecmd.mk
