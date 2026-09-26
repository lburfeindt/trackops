TOOLS := tools/format-playlist tools/id3-sanitizer

.PHONY: test verify format lint build

test verify format lint build:
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" $@; done
