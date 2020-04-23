.DEFAULT_GOAL := all

EXECUTABLE=pockety
WINDOWS=./bin/windows_amd64
LINUX=./bin/linux_amd64
DARWIN=./bin/darwin_amd64
VERSION=$(shell git describe --tags --abbrev=0)

prepare:
	@echo Cleaning the bin directory
	@rm -rfv ./bin/*

windows:
	@echo Building Windows amd64 binaries
	@env GOOS=windows GOARCH=amd64 go build -i -v -o $(WINDOWS)/$(EXECUTABLE).exe -ldflags="-s -w"  *.go

linux:
	@echo Building Linux amd64 binaries
	@env GOOS=linux GOARCH=amd64 go build -i -v -o $(LINUX)/$(EXECUTABLE) -ldflags="-s -w"  *.go

darwin:
	@echo Building Mac amd64 binaries
	@env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN)/$(EXECUTABLE) -ldflags="-s -w"  *.go

build: ## Builds the binaries.
build: windows linux darwin

package:
	@echo Creating the zip file
	@tar -C $(DARWIN) -cvzf ./bin/$(EXECUTABLE)_darwin-$(VERSION).tar.gz $(EXECUTABLE)
	@zip -j ./bin/$(EXECUTABLE)_windows-$(VERSION).zip $(WINDOWS)/$(EXECUTABLE).exe
	@tar -C $(LINUX) -cvzf ./bin/$(EXECUTABLE)_linux-$(VERSION).tar.gz $(EXECUTABLE)

install:
	@cp -pv $(DARWIN)/$(EXECUTABLE)

help: ##  Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

all: prepare build package clean

clean: ## Removes the artifacts.
	@rm -rf $(WINDOWS) $(LINUX) $(DARWIN)

.PHONY: all
