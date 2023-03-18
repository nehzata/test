SRC = $(shell find . -type f -name '*.go' -not -path "./databases/*")
REGEN_DIRS = $(shell find . -name '*.go' | xargs grep -l //go:generate | xargs dirname | sort | uniq)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: 
dev: ## Run hot reload dev server
	reflex -d fancy -c reflex.conf

.PHONY: lint
lint: ## Lint the code.
	golangci-lint run

.PHONY:
fmt: ## Format all code.
	@gofmt -l -w $(SRC)