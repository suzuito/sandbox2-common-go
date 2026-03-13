
tools/httpfakeserver/dist/prd/hfs: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/httpfakeserver/cmd/$(notdir $@)/*.go

tools/httpfakeserver/dist/e2e/hfs: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/httpfakeserver/cmd/$(notdir $@)/*.go

.PHONY: e2e-httpfakeserver-hfs
e2e-httpfakeserver-hfs: tools/httpfakeserver/dist/e2e/hfs
	FILE_PATH_SERVER_BIN=$(abspath $<) \
	PORT=8100 \
sh run-e2e.sh $< tools/httpfakeserver/cov/e2e/hfs ./e2e/httpfakeserver/hfs/...
