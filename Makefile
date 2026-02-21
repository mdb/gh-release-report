SOURCE=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')
VERSION?=0.0.2
NAME=gh-release-report
GORELEASER=go run github.com/goreleaser/goreleaser/v2@v2.13.3

default: build

version:
	@echo $(VERSION)
.PHONY: version

build:
	$(GORELEASER) release \
		--snapshot \
		--skip=publish \
		--clean
.PHONY: build

test: vet fmtcheck
	go test -v -coverprofile=coverage.out -count=1 $(SOURCE)
.PHONY: test

acc-test:
	go test -v --tags=acceptance -count=1 ./cmd
.PHONY: acc-test

vet:
	go vet $(SOURCE)
.PHONY: vet

fmt:
	gofmt -w $(GOFMT_FILES)
.PHONY: fmt

fmtcheck:
	test -z $(shell go fmt $(SOURCE))
.PHONY: fmtcheck

check-tag:
	./scripts/ensure-unique-version.sh "$(VERSION)"
.PHONY: check-tag

tag: check-tag
	echo "creating git tag $(VERSION)"
	git tag $(VERSION)
	git push origin $(VERSION)
.PHONY: tag

release:
	$(GORELEASER) release \
		--clean
.PHONY: release

install:
	mkdir -p ~/.local/share/gh/extensions/$(NAME)
	cp dist/$(NAME)_$(shell echo $(shell uname) | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')*/$(NAME) ~/.local/share/gh/extensions/$(NAME)/
.PHONY: install

demo:
	vhs < demo.tape
.PHONY: demo
