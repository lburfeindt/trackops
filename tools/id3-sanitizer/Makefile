.PHONY: test verify format lint build clean

format:
	gofmt -w .

lint:
	go vet ./...

test:
	go test ./...

verify:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then printf '%s\n' "$$unformatted"; exit 1; fi
	$(MAKE) lint
	$(MAKE) test

build:
	mkdir -p bin
	go build -o bin/id3-sanitizer .

clean:
	go clean -cache -testcache -modcache
	rm -rf bin
