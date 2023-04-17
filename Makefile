.PHONY: build all release

GIT_HASH := $(shell git rev-parse HEAD)
NEAREST_VERSION := $(shell git tag --contains $(GIT_HASH) | tail -n 1)
all: build

release: .build/envctl-$(NEAREST_VERSION).tgz 

.build/README.md: README.md
	cp README.md .build/README.md

.build/envctl-$(NEAREST_VERSION).tgz: .build/bin/envctl .build/README.md
	tar -czf .build/envctl-$(NEAREST_VERSION).tgz -C .build bin README.md

build: .build/bin/envctl

.build/bin/envctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@