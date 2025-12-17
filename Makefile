.PHONY: help
help: ## Show this help
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: test
test:
	@xvfb-run go test ./...

.PHONY: pprof
pprof: ## Runs pprof server for 'cpu.out'
	@go tool pprof --http=:8081 cpu.out

.PHONY: lint
lint: ## Runs the linter
	@go tool staticcheck ./...
