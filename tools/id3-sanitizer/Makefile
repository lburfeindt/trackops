.PHONY: test verify format lint build clean package-macos

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
	go build -o bin/id3-sanitizer-gui ./cmd/id3-sanitizer-gui

package-macos:
	@test "$$(uname -s)" = Darwin || { echo 'package-macos requires macOS'; exit 1; }
	mkdir -p "bin/ID3 Sanitizer.app/Contents/MacOS"
	cp packaging/macos/Info.plist "bin/ID3 Sanitizer.app/Contents/Info.plist"
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o "bin/ID3 Sanitizer.app/Contents/MacOS/id3-sanitizer-gui" ./cmd/id3-sanitizer-gui
	codesign --force --sign - "bin/ID3 Sanitizer.app"

clean:
	go clean -cache -testcache -modcache
	rm -rf bin
