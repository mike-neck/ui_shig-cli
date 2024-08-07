SHELL=/usr/bin/env bash

SCRIPTS := $(wildcard scripts/*.sh)
NAMES := $(patsubst scripts/%,%,$(patsubst %.sh,%,$(SCRIPTS)))

help:
	@echo $(@)
	@echo "$(NAMES)"

.PHONY: $(NAMES)
$(NAMES):
	@echo $(@)
	@"./scripts/$(@).sh"

####################################
# Test
TEST_FILES := $(wildcard *_test.go)
TEST_NAMES := $(patsubst %_test.go,run-test-%,$(TEST_FILES))
PRODUCT_FILES := $(filter-out $(TEST_FILES),$(wildcard *.go))

run-test-%: %_test.go
	@echo "$(@)"
	@go test "$(patsubst run-test-%,%_test.go,$(@))" $(PRODUCT_FILES)

.PHONY: all-test
all-test: $(TEST_NAMES)

.PHONY: test
test: all-test

compile-test-%:%_test.go
	@echo "$(@) -> $(PWD)/bin/tests/$(patsubst compile-test-%,%,$(@))"
	@GOOS=windows GOARCH=amd64 go test -c -o "$(PWD)/bin/tests/$(patsubst compile-test-%,%,$(@)).exe" "$(patsubst compile-test-%,%_test.go,$(@))" $(PRODUCT_FILES)

.PHONY: compile-all-tests
compile-all-tests: $(patsubst %_test.go,compile-test-%,$(TEST_FILES))

####################################
# OS ARCH 別タスク

TARGET_OS := windows darwin linux
TARGET_ARCH := arm64 amd64
CURRENT_OS := $(shell go env GOOS)
CURRENT_ARCH := $(shell go env GOARCH)

#リナックスは別途ビルド手順が必要なので if 分岐してる(TODO Linux用ビルドの手順(oto 依存))
define BuildWithOsArch
.PHONY: build-$(1)-$(2)
build-$(1)-$(2):
	@printf "%-20s " $$(@)
	@if [[ "$(1)" != "linux" ]] || [[ "$(1)" == "$(CURRENT_OS)" && "$(2)" == "$(CURRENT_ARCH)" ]]; then \
		scripts/build.sh "$(1)" "$(2)" ;\
		head -n 20 < "README.md" > "bin/$(1)/$(2)/README.md" ;\
		echo "...done"; \
	else \
		echo "...skip"; \
	fi

endef

define CleanWithOsArch
.PHONY: clean-$(1)-$(2)
clean-$(1)-$(2):
	@echo $$(@)
	@rm -rf "bin/$(1)/$(2)"

endef

define ArchiveWithOsArch
.PHONY: archive-$(1)-$(2)
archive-$(1)-$(2):
	@echo $$(@)
	@"./scripts/archive.sh" $(1) $(2)

endef

define UploadWithOsArch
.PHONY: upload-$(1)-$(2)
upload-$(1)-$(2):
	@echo $$(@)
	@"./scripts/upload.sh" $(1) $(2)

endef

TEMPLATES := BuildWithOsArch CleanWithOsArch ArchiveWithOsArch UploadWithOsArch
$(foreach template,$(TEMPLATES),$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(eval $(call $(template),$(os),$(arch))))))

############################
# all tasks

ALL_TASKS := build clean archive upload

define AllTasks
.PHONY: $(1)-all
$(1)-all: $(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(1)-$(os)-$(arch)))

endef

$(foreach task,$(ALL_TASKS),$(eval $(call AllTasks,$(task))))

.PHONY: show-build-result
show-build-result:
	@echo "$(@)"
	@find "$(PWD)/bin" -type f | sort

.PHONY: tag
tag:
	@echo $(@)
	@$(PWD)/scripts/create-new-tag.sh
