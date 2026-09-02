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
	git remote | grep -q '^dev-charter$$' || \
	  git remote add dev-charter https://github.com/y-marui/dev-charter
	git fetch dev-charter
	@STASHED=0; \
	if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		git stash push -u -m "update-charter"; \
		STASHED=1; \
	fi; \
	git subtree pull --prefix=docs/dev-charter dev-charter main --squash; \
	if [ "$$STASHED" = "1" ]; then git stash pop; fi

# alfred-workflow-notes lives at docs/alfred-workflow-notes/ *inside* the
# alfred-workflow-template repo, not at that repo's root (unlike
# dev-charter, whose repo root IS the shared content) — a plain
# `git subtree pull` would pull the whole template repo in. Split that
# subdirectory's history out into a throwaway local branch first, then
# merge just that.
update-workflow-notes:
	git remote | grep -q '^alfred-workflow-notes$$' || \
	  git remote add alfred-workflow-notes https://github.com/y-marui/alfred-workflow-template
	git fetch alfred-workflow-notes
	@STASHED=0; \
	if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		git stash push -u -m "update-workflow-notes"; \
		STASHED=1; \
	fi; \
	git subtree split --prefix=docs/alfred-workflow-notes --branch workflow-notes-split alfred-workflow-notes/main; \
	git subtree merge --prefix=docs/alfred-workflow-notes workflow-notes-split --squash; \
	git branch -D workflow-notes-split; \
	if [ "$$STASHED" = "1" ]; then git stash pop; fi
