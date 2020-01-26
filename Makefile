# Copyright 2019 Manlio Perillo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# A Makefile template for Go projects.

# Exported variable definitions.
export GO111MODULE := on

# Imported variables.
# GOPKG - used to select the target package

# Variable definitions.
BENCHFLAGS := -v
COVERMODE := atomic # atomic is necessary if the -race flag is enabled
TESTFLAGS := -race -v

# Standard rules.
.POSIX:

.PHONY: build bench clean cover github install lint print test test-trace trace vet

# Default rule.
build:
	go build -o build ./...

# Custom rules.
bench:
	go test ${BENCHFLAGS} -bench=. -benchmem ./...

clean:
	go mod tidy
	go clean
	go clean -i
	rm -f build/*

cover:
	go tool cover -html=build/coverage.out -o=build/coverage.html

github:
	git push --follow-tags -u github master

install:
	go install ./...

lint:
	golint ./...

print:
	goprint -font='"Inconsolata" 10pt/12pt' ${GOPKG} > build/pkg.html
	prince -o build/pkg.pdf build/pkg.html

test:
	go test ${TESTFLAGS} -covermode=${COVERMODE} \
		-coverprofile=build/coverage.out ./...

test-trace:
	go test ${TESTFLAGS} -trace=build/trace.out ${GOPKG}

trace:
	go tool trace build/trace.out

vet:
	go vet ./...
