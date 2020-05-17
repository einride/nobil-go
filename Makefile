SHELL := /bin/bash

.PHONY: all
all: \
	go-lint \
	go-test \
	go-mod-tidy \

.PHONY: clean
clean:
	rm -rf tools/*/*/

include tools/golangci-lint/rules.mk

.PHONY: go-lint
go-lint: $(golangci_lint)
	GOGC=1 GOMAXPROCS=1 $(golangci_lint) run --timeout 5m

.PHONY: go-test
go-test:
	go test -race ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v

.PHONY: gcloud-builds-triggers-create
gcloud-builds-triggers-create: repo_name = $(shell basename -s .git $(shell git config --get remote.origin.url))
gcloud-builds-triggers-create:
	gcloud beta builds triggers create github \
		--project=einride \
		--repo-owner='einride' \
		--repo-name='$(repo_name)' \
		--pull-request-pattern='.*' \
		--description='$(repo_name)-review' \
		--build-config='.cloudbuild/review.yaml'
