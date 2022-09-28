#!/usr/bin/make -f

.ONESHELL:
.SHELL := /usr/bin/bash

AUTHOR := "noelruault"
PROJECTNAME := $(shell basename "$$(pwd)")
PROJECTPATH := $(shell pwd)

help:
	@echo "Usage: make [options] [arguments]\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

example: ## Run a example using the example image
	@mkdir -p $(PROJECTPATH)/tmp
	@go run main.go -image=./assets/starry_night-vincent_van-gogh.png -xlen=320 -ylen=253 -out=./tmp/

release: ## Tags to trigger a new release
	@read -p "Release version: " VERSION;\
	git tag $$VERSION && git push origin --tags
