TOOLS := tools/format-playlist tools/id3-sanitizer

.PHONY: test verify format lint build clean

test verify format lint build clean:
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" $@; done
