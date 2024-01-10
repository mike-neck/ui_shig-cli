SCRIPTS := $(wildcard scripts/*.sh)
NAMES := $(patsubst scripts/%,%,$(patsubst %.sh,%,$(SCRIPTS)))

help:
	@echo $(@)
	@echo "$(NAMES)"

.PHONY: $(NAMES)
$(NAMES):
	@echo $(@)
	@"./scripts/$(@).sh"
