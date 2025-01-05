
tools/fakecmd/dist/prd/fakecmd: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/fakecmd/cmd/$(notdir $@)/*.go

tools/fakecmd/dist/e2e/fakecmd: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/fakecmd/cmd/$(notdir $@)/*.go

.PHONY: e2e-fakecmd-fakecmd
e2e-fakecmd-fakecmd:
	make start-e2e-environment
	sh run-e2e.sh tools/fakecmd/dist/e2e/fakecmd tools/fakecmd/cov/e2e/fakecmd ./e2e/fakecmd/fakecmd/...
