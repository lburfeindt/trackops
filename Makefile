.PHONY: test check build

test:
	go test ./...

check:
	@unformatted="$$(gofmt -l $$(find . -name '*.go' -not -path './.git/*'))"; \
	if [ -n "$$unformatted" ]; then printf '%s\n' "$$unformatted"; exit 1; fi
	go vet ./...
	go test ./...

build:
	mkdir -p bin
	go build -o bin/format-playlist ./cmd/format-playlist
