
tools/release/dist/prd/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/release/cmd/$(notdir $@)/*.go

tools/release/dist/e2e/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/release/cmd/$(notdir $@)/*.go

.PHONY: e2e-release-increment-release-version
e2e-release-increment-release-version:
	make tools/release/dist/e2e/increment-release-version
	mkdir -p tools/release/cov/e2e/increment-release-version
	make start-e2e-environment
	sh run-e2e.sh tools/release/dist/e2e/increment-release-version tools/release/cov/e2e/increment-release-version
