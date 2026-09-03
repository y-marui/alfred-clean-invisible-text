.PHONY: build test lint fmt fetch-cli build-workflow precommit update-charter update-workflow-notes

build:
	go build ./...

test:
	go test -p 1 ./...

lint:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to be run on:"; \
		echo "$$out"; \
		exit 1; \
	fi
	go vet ./...

fmt:
	gofmt -w .

fetch-cli:
	scripts/fetch-cli-binaries.sh

build-workflow:
	scripts/build-workflow.sh

precommit:
	pre-commit run --all-files

update-charter:
	curl -fsSL https://raw.githubusercontent.com/y-marui/dev-charter/main/scripts/install.sh | CHARTER_UPDATE_ONLY=1 bash

update-workflow-notes:
	curl -fsSL https://raw.githubusercontent.com/y-marui/alfred-workflow-template/main/scripts/install-workflow-notes.sh | bash
