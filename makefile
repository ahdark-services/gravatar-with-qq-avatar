GO=go

.PHONY: build

build-%:
	$(GO) build -o bin/$* cmd/$*/main.go

run-%:
	$(GO) run cmd/$*/main.go

build:
	@for dir in $(shell ls cmd); do \
		$(GO) build -o bin/$$dir cmd/$$dir/main.go; \
	done
