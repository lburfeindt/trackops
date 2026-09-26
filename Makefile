TOOLS := tools/format-playlist

.PHONY: test check build

test check build:
	@set -e; for tool in $(TOOLS); do $(MAKE) -C "$$tool" $@; done
