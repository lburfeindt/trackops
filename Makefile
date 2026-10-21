WAILS ?= go run github.com/wailsapp/wails/v2/cmd/wails@v2.17.0
.PHONY: test verify format lint prepare-build build package-macos clean test-ui

format:
	gofmt -w cmd/trackops frontend/assets.go internal/desktop internal/id3sanitizer internal/playlist

lint:
	go vet ./...

test:
	go test ./...

verify:
	@test -z "$$(gofmt -l cmd/trackops frontend/assets.go internal/desktop internal/id3sanitizer internal/playlist)"
	$(MAKE) lint
	$(MAKE) test

prepare-build:
	mkdir -p build/darwin
	cp packaging/macos/appicon.png build/appicon.png
	cp packaging/macos/Info.plist build/darwin/Info.plist

build: prepare-build
	cd cmd/trackops && CGO_CFLAGS='-mmacosx-version-min=13.0' CGO_LDFLAGS='-framework UniformTypeIdentifiers -mmacosx-version-min=13.0' $(WAILS) build

package-macos: prepare-build
	@test "$$(uname -s)" = Darwin || { echo 'package-macos requires macOS'; exit 1; }
	cd cmd/trackops && CGO_CFLAGS='-mmacosx-version-min=13.0' CGO_LDFLAGS='-framework UniformTypeIdentifiers -mmacosx-version-min=13.0' $(WAILS) build -clean -platform darwin/arm64
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
	rm -rf build
