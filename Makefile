.PHONY: build all

all: build

build: .build/bin/envctl

.build/bin/envctl: main.go .build/bin
	go build -o $@ main.go

.build:
	mkdir -p $@

.build/bin: .build
	mkdir -p $@