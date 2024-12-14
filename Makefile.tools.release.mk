
tools/release/dist/prd/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/release/cmd/$(notdir $@)/*.go

tools/release/dist/e2e/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/release/cmd/$(notdir $@)/*.go

.PHONY: e2e-release-increment-release-version
e2e-release-increment-release-version:
	make tools/release/dist/e2e/increment-release-version
	make $(BIN_E2E_INCREMENT_RELEASE_VERSION)
	make start-e2e-environment
	export FILE_PATH_BIN=$(abspath tools/release/dist/e2e/increment-release-version) && \
	export GOCOVERDIR=$(abspath tools/release/cov/e2e/increment-release-version) && \
	mkdir -p $${GOCOVERDIR} && rm $${GOCOVERDIR}/* && \
	go test -count=1 -v ./tools/release/e2e/increment-release-version/... || : && \
	sh report-e2e.sh $${GOCOVERDIR}
