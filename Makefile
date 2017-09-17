all: remarked

remarked: $(shell find . -name '*.go')
	cd cmd/remarked && go build -o ../../remarked

install:
	cd cmd/remarked && go install

test:
	go test -v ./...

.PHONY: test
.PHONY: clean
.PHONY: install
.PHONY: all
