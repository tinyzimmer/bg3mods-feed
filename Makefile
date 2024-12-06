build:
	go build -o bin/$(shell basename $(PWD)) main.go

PARALLELISM ?= 4
.PHONY: dist
dist:
	goreleaser release --clean --parallelism=$(PARALLELISM) --snapshot