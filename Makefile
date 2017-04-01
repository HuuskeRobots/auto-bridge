PROJECT := auto-bridge
SCRIPTDIR := $(shell pwd)
ROOTDIR := $(shell cd $(SCRIPTDIR) && pwd)
VERSION:= $(shell cat $(ROOTDIR)/VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)
BINDIR := $(ROOTDIR)
VENDORDIR := $(SCRIPTDIR)/vendor

ORGPATH := github.com/HuuskeRobots
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)
BIN := $(BINDIR)/$(PROJECT)

GOPATH := $(GOBUILDDIR)
GOVERSION := 1.8-alpine

ifndef GOOS
	GOOS := darwin
endif
ifndef GOARCH
	GOARCH := amd64
endif
ifeq ("$(GOOS)", "windows")
        GOEXE := .exe
endif

SOURCES := $(shell find $(SRCDIR) -name '*.go')

.PHONY: all clean deps sample

all: deployment

build: $(BIN)

deployment:
	@${MAKE} -B GOOS=windows GOARCH=amd64 GOEXE=.exe build
	@${MAKE} -B GOOS=linux GOARCH=amd64 build
	@${MAKE} -B GOOS=darwin GOARCH=amd64 build

clean:
	rm -Rf $(BIN) $(GOBUILDDIR)

deps:
	@${MAKE} -B -s $(GOBUILDDIR)

$(GOBUILDDIR):
	@mkdir -p $(ORGDIR)
	@rm -f $(REPODIR) && ln -s ../../../.. $(REPODIR)

update-vendor:
	@rm -Rf $(VENDORDIR)
	@pulsar go vendor -V $(VENDORDIR) \
		gobot.io/x/gobot/ \
		gobot.io/x/gobot/platforms/firmata \
		golang.org/x/sys/unix

$(BIN): $(GOBUILDDIR) $(SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOPATH=$(GOBUILDDIR) go build -ldflags=-s -o bin/$(GOOS)-$(GOARCH)/auto$(GOEXE) $(REPOPATH)

