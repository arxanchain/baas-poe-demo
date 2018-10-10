#
# Copyright arxanfintech.com. 2017 All Rights Reserved.
#
# Purpose: make file
#
# This file is subject to the terms and conditions defined in
# file 'LICENSE.txt', which is part of this source code package.
#

PROJECT_NAME=baas-poe-demo
PKGNAME = github.com/arxanchain/$(PROJECT_NAME)

EXECUTABLES = go git
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH: Check dependencies")))

all: baas-poe-demo

.PHONY: baas-poe-demo
baas-poe-demo:
	GOBIN=$(abspath .) go install $(PKGNAME)
	@echo "Binary available $(PROJECT_NAME)"

.PHONY: clean
clean:
	-@rm -f $(PROJECT_NAME)

