.PHONY: build all

all: build

release: .build/bin/envctl
	tar -czf .build/bin/envctl.tgz -C .build/bin envctl

build: .build/bin/envctl

.build/bin/envctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@