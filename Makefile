SHELL := /bin/bash

.PHONY: all
all: \
	commitlint \
	prettier-markdown \
	go-lint \
	go-review \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

.PHONY: clean
clean:
	rm -rf tools/*/*/

include tools/commitlint/rules.mk
include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/prettier/rules.mk
include tools/semantic-release/rules.mk

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@go test -cover -race ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy -v
