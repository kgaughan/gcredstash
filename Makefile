NAME:=gcredstash

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/*/*.go)
DOCS:=$(wildcard docs/*.md mkdocs.yml)

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: go.mod $(NAME) ## Build the gcredstash binary

.PHONY: tidy
tidy: go.mod fmt ## Tidy go.mod and format the code

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(NAME) dist site .venv coverage.out coverage.html

$(NAME): $(SOURCE) go.sum
	CGO_ENABLED=0 go build -v -tags netgo -trimpath -ldflags '-s -w' -o $@ ./cmd/$@

.PHONY: update
update: ## Update dependencies
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify
	@touch go.sum

go.mod: $(SOURCE)
	go mod tidy

.PHONY: fmt
fmt: ## Format the code
	go fmt ./...

.PHONY: lint
lint: ## Lint the code
	go vet ./...
	golangci-lint run ./...

.PHONY: serve-docs
serve-docs: .venv ## Serve the documentation locally
	uv run mkdocs serve

.PHONY: docs
docs: .venv $(DOCS) ## Build the documentation site
	uv run mkdocs build

.venv: requirements.txt
	uv venv
	uv pip install -r requirements.txt

%.txt: %.in
	uv pip compile $< -o $@

.PHONY: tests
tests: ## Run the tests
	go test -cover -coverprofile=coverage.out -v ./...

coverage.out: tests

.PHONY: coverage-html
coverage-html: coverage.out ## Generate HTML report from coverage data
	sed -E '/\/(mockaws|testutils)\//d' coverage.out | go tool cover -html=/dev/stdin -o coverage.html

.PHONY: test-release
test-release: ## Run `goreleaser release` without publishing anything
	goreleaser release --auto-snapshot --clean --skip docker --skip publish

.PHONY: mocks
mocks: ## Generate mock implementations
	go install go.uber.org/mock/mockgen@v0.6.0
	mockgen -package mockaws -destination internal/mockaws/dynamodbmock.go github.com/kgaughan/gcredstash/internal DynamoDB
	mockgen -package mockaws -destination internal/mockaws/kmsmock.go github.com/kgaughan/gcredstash/internal Kms
