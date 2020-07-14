.DEFAULT_GOAL	:= build
ORG = github.com/staroids
PROJECT = starctl
REPOPATH ?= $(ORG)/$(PROJECT)
BUILD_PACKAGE = $(REPOPATH)/cmd/starctl

.PHONY: build
build:
	@go build -o starctl $(BUILD_PACKAGE)

.PHONY: buildlinux
buildlinux:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o starctl $(BUILD_PACKAGE)

.PHONY: clean
clean:
	@rm -f starctl starctl-*
	@go clean -modcache

.PHONY: test
test:
	@go test $(BUILD_PACKAGE) -cover

.PHONY: release
release:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o starctl-linux-amd64 $(BUILD_PACKAGE)
	@env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o starctl-darwin-amd64 $(BUILD_PACKAGE)
