NAME:=gcredstash

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/*/*.go)
DOCS:=$(wildcard docs/*.md mkdocs.yml)

.PHONY: build
build: go.mod $(NAME)

.PHONY: tidy
tidy: go.mod fmt

.PHONY: clean
clean:
	rm -rf $(NAME) dist site

$(NAME): $(SOURCE) go.sum
	CGO_ENABLED=0 go build -v -tags netgo -trimpath -ldflags '-s -w' -o $@ ./cmd/$@

.PHONY: update
update:
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify
	@touch go.sum

go.mod: $(SOURCE)
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: serve-docs
serve-docs: .venv
	.venv/bin/mkdocs serve

.PHONY: docs
docs: .venv $(DOCS)
	.venv/bin/mkdocs build

.venv: requirements.txt
	uv venv
	uv pip install -r requirements.txt

%.txt: %.in
	uv pip compile $< > $@

.PHONY: tests
tests:
	go test -cover -v ./...

.PHONY: mocks
mocks:
	go install go.uber.org/mock/mockgen@v0.6.0
	mockgen -package mockaws -destination internal/mockaws/dynamodbmock.go github.com/kgaughan/gcredstash/internal DynamoDB
	mockgen -package mockaws -destination internal/mockaws/kmsmock.go github.com/kgaughan/gcredstash/internal Kms
