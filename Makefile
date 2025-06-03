BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint

.PHONY: mac-init
mac-init:

# golangci-lint だけは go.mod で管理しない。`go install` によるインストールが非推奨とされているため
# https://golangci-lint.run/welcome/install/#install-from-sources
$(BIN_GOLANGCI_LINT): Makefile
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

.PHONY: lint
lint: $(BIN_GOLANGCI_LINT)
	$(BIN_GOLANGCI_LINT) run ./...

.PHONY: godoc
godoc:
	go tool air -c air.godoc.toml

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

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

.PHONY: build
build:
	rm -rf dist/prd/$(GOOS)/$(GOARCH)
	mkdir -p dist/prd/$(GOOS)/$(GOARCH)
	make tools/fakecmd/dist/prd/fakecmd
	mv tools/fakecmd/dist/prd/fakecmd dist/prd/$(GOOS)/$(GOARCH)/
	make tools/release/dist/prd/increment-release-version
	mv tools/release/dist/prd/increment-release-version dist/prd/$(GOOS)/$(GOARCH)/
	make tools/terraform/dist/prd/check-terraform-rules
	mv tools/terraform/dist/prd/check-terraform-rules dist/prd/$(GOOS)/$(GOARCH)/
	make tools/terraform/dist/prd/terraform_on_github_action
	mv tools/terraform/dist/prd/terraform_on_github_action dist/prd/$(GOOS)/$(GOARCH)/
