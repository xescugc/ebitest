.PHONY: help
help: ## Show this help
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: test
test:
	@xvfb-run go test ./... -count=1

.PHONY: pprof
pprof:
	@go tool pprof --http=:8081 cpu.out
