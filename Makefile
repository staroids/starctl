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
	@rm -f starctl
	@go clean -modcache

.PHONY: test
test:
	@go test $(BUILD_PACKAGE) -cover
