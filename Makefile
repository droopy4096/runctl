.PHONY: build all release check_version

GIT_HASH := $(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --contains $(GIT_HASH) | tail -n 1)

all: build

# check_version:
# 	ifeq ($(VERSION),)
# 		$(error No version specified)
# 	else
# 		echo "Version: $(VERSION)"
# 	endif

release: .build/runctl-$(VERSION).tgz 

.build/README.md: README.md
	cp README.md .build/README.md

.build/runctl-$(VERSION).tgz: .build/bin/runctl .build/README.md
	tar -czf .build/runctl-$(VERSION).tgz -C .build bin README.md

build: .build/bin/runctl

.build/bin/runctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@