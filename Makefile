.PHONY: build all

GIT_HASH := $(shell git rev-parse HEAD)
NEAREST_VERSION := $(shell git tag --contains $(GIT_HASH))
all: build

release: .build/bin/envctl
	cp README.md .build/README.md
	tar -czf .build/envctl-$(NEAREST_VERSION).tgz -C .build bin README.md

build: .build/bin/envctl

.build/bin/envctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@