.PHONY: build
build:
	go build -o ./bin/computer-club ./cmd/main.go

TEST_DIR := examples
.PHONY: test
test: build
	@for f in $(shell ls ./${TEST_DIR}); \
	do \
	echo "Running test from file $${f}:"; ./bin/computer-club ./${TEST_DIR}/$${f}; echo ""; echo ""; \
	done

.PHONY: unit-test
unit-test:
	go test ./...


version=
.PHONY: docker-build
docker-build:
	docker build -t computer-club-system .
ifdef version
	docker image tag computer-club-system:latest computer-club-system:$(version)
endif

file:=
.PHONY: docker-run
docker-run:
	docker run -v "$(dir $(realpath $(lastword $(MAKEFILE_LIST))))${TEST_DIR}:/${TEST_DIR}" computer-club-system "/${TEST_DIR}/${file}"
