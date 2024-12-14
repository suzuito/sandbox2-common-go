
dist/prd/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/release/cmd/increment-release-version/*.go

dist/e2e/increment-release-version: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/release/cmd/increment-release-version/*.go

.PHONY: e2e-increment-release-version
e2e-increment-release-version: dist/e2e/increment-release-version
	make start-e2e-environment
	P=cov/e2e/increment-release-version && rm -r $${P} && mkdir -p $${P}
	FILE_PATH_BIN=$(abspath dist/e2e/increment-release-version) \
	GOCOVERDIR=$(abspath cov/e2e/increment-release-version) \
	go test -count=1 -v ./tools/release/e2e/increment-release-version/... || :
	sh report-e2e.sh cov/e2e/increment-release-version
