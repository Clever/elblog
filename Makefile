include golang.mk
.DEFAULT_GOAL := test

SHELL := /bin/bash
PKG := github.com/Clever/elblog
PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test $(PKGS) run install_deps

$(eval $(call golang-version-check,1.12))

test: $(PKGS)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)
