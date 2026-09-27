.PHONY: test verify format lint build clear

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
	go build -o bin/format-playlist .

clear:
	go clean -cache -testcache -modcache
	rm -rf bin
