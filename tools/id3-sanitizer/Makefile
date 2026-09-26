.PHONY: test verify format lint build

GO_FILES := $(shell find . -name '*.go' -not -path './.git/*')

format:
	gofmt -w $(GO_FILES)

lint:
	go vet ./...

test:
	go test ./...

verify:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then printf '%s\n' "$$unformatted"; exit 1; fi
	$(MAKE) lint
	$(MAKE) test

build:
	mkdir -p bin
	go build -o bin/id3-sanitizer .
