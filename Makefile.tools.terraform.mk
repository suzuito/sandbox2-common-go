
tools/terraform/dist/prd/check-terraform-rules: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ tools/terraform/cmd/$(notdir $@)/*.go

tools/terraform/dist/e2e/check-terraform-rules: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ tools/terraform/cmd/$(notdir $@)/*.go

.PHONY: e2e-terraform-check-terraform-rules
e2e-terraform-check-terraform-rules:
	make start-e2e-environment
	sh run-e2e.sh tools/terraform/dist/e2e/check-terraform-rules tools/terraform/cov/e2e/check-terraform-rules ./e2e/terraform/check-terraform-rules/...
