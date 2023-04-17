.PHONY: build all release check_version

GIT_HASH := $(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --contains $(GIT_HASH) | tail -n 1)

check_version:
ifeq ($(VERSION),)
	$(error No version specified)
else
	echo "Version: $(VERSION)"
endif

all: build

release: .build/envctl-$(VERSION).tgz 

.build/README.md: README.md
	cp README.md .build/README.md

.build/envctl-$(VERSION).tgz: check_version .build/bin/envctl .build/README.md
	tar -czf .build/envctl-$(VERSION).tgz -C .build bin README.md

build: .build/bin/envctl

.build/bin/envctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@