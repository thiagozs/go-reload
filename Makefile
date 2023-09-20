PROJECT_NAME := reload-runner
INSTALL_PATH := ~/.local/bin
GO_CMD := go
INSTALL_CMD := install
RM := rm -f
GO_LINT := golangci-lint run
GO_BUILD := $(GO_CMD) build
GO_INSTALL := $(GO_CMD) install
GO_CLEAN := $(GO_CMD) clean
GO_TEST := $(GO_CMD) test
GO_FMT := $(GO_CMD) fmt
GO_BUILD_TARET := build/$(PROJECT_NAME)

all: build

build:
	$(GO_BUILD) -o $(GO_BUILD_TARET)

install:
	$(GO_INSTALL)
	$(INSTALL_CMD) -m 755 $(GO_BUILD_TARET) $(INSTALL_PATH)

test:
	$(GO_TEST) -v ./...

fmt:
	$(GO_FMT) ./...

lint:
	$(GO_LINT)

clean:
	$(GO_CLEAN)
	$(RM) $(GO_BUILD_TARET)

.PHONY: all build install test fmt lint clean
