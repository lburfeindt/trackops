TOOLS := tools/format-playlist tools/id3-sanitizer
WAILS ?= wails

.PHONY: test verify format lint build clean test-ui

format:
	gofmt -w main.go app.go internal/playlist
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" format; done

lint:
	go vet ./...
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" lint; done

test:
	go test ./...
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" test; done

verify: lint test

build:
	CGO_CFLAGS='-mmacosx-version-min=13.0' CGO_LDFLAGS='-framework UniformTypeIdentifiers -mmacosx-version-min=13.0' $(WAILS) build
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" build; done

test-ui: build
	./test/ui/open-playlist.sh

clean:
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" $@; done
	rm -rf build/bin
