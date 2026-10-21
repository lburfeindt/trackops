WAILS ?= go run github.com/wailsapp/wails/v2/cmd/wails@v2.17.0
.PHONY: test verify format lint build package-macos clean test-ui

format:
	gofmt -w main.go app.go mp3_picker_*.go internal/id3sanitizer internal/playlist

lint:
	go vet ./...

test:
	go test ./...

verify:
	@test -z "$$(gofmt -l main.go app.go mp3_picker_*.go internal/id3sanitizer internal/playlist)"
	$(MAKE) lint
	$(MAKE) test

build:
	CGO_CFLAGS='-mmacosx-version-min=13.0' CGO_LDFLAGS='-framework UniformTypeIdentifiers -mmacosx-version-min=13.0' $(WAILS) build

package-macos:
	@test "$$(uname -s)" = Darwin || { echo 'package-macos requires macOS'; exit 1; }
	CGO_CFLAGS='-mmacosx-version-min=13.0' CGO_LDFLAGS='-framework UniformTypeIdentifiers -mmacosx-version-min=13.0' $(WAILS) build -clean -platform darwin/arm64
	codesign --force --deep --sign - build/bin/TrackOps.app
	codesign --verify --deep --strict build/bin/TrackOps.app
	@test "$$(/usr/libexec/PlistBuddy -c 'Print CFBundleIdentifier' build/bin/TrackOps.app/Contents/Info.plist)" = com.lburfeindt.trackops
	@test "$$(/usr/libexec/PlistBuddy -c 'Print CFBundleShortVersionString' build/bin/TrackOps.app/Contents/Info.plist)" = 0.1.0
	@test "$$(/usr/libexec/PlistBuddy -c 'Print LSMinimumSystemVersion' build/bin/TrackOps.app/Contents/Info.plist)" = 13.0
	@test "$$(lipo -archs build/bin/TrackOps.app/Contents/MacOS/TrackOps)" = arm64

test-ui: package-macos
	./test/ui/open-playlist.sh
	./test/ui/playlist-states.sh
	./test/ui/save-playlist.sh
	./test/ui/sanitize-mp3.sh
	./test/ui/sanitize-folder.sh

clean:
	rm -rf build/bin
