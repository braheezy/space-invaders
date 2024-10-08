PACKAGE := space-invaders

# Go defintions
GOCMD ?= go
GOBUILD := $(GOCMD) build
GOINSTALL := $(GOCMD) install
ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif
GOARCH := amd64
DLV_PORT := 44572

# Build definitions
BUILD_ENTRY := $(PWD)
BIN_DIR := $(PWD)/bin

# Determine the file extension based on the platform
ifeq ($(OS),Windows_NT)
  EXTENSION := .exe
else
  EXTENSION :=
endif
# Different platform support
PLATFORMS := linux windows darwin
BINARIES := $(addprefix $(BIN_DIR)/,$(addsuffix /$(PACKAGE)$(EXTENSION),$(PLATFORMS)))

# Fancy colors
BOLD := $(shell tput bold)
ITALIC := \e[3m
YELLOW := $(shell tput setaf 222)
GREEN := $(shell tput setaf 114)
BLUE := $(shell tput setaf 111)
PURPLE := $(shell tput setaf 183)
END := $(shell tput sgr0)

# Function to colorize a command help string
command-style = $(BLUE)$(BOLD)$1$(END)  $(ITALIC)$(YELLOW)$2$(END)

define help_text
$(PURPLE)$(BOLD)Targets:$(END)
  - $(call command-style,all,      Build $(PACKAGE) for all targets (Linux, Windows, Mac, 64-bit))
  - $(call command-style,build,    Build $(PACKAGE) for current host architecture)
  - $(call command-style,run,      Build and run $(PACKAGE) for current host)
  - $(call command-style,install,  Build and install $(PACKAGE) for current host)
  - $(call command-style,debug,    Run a dlv debug headless session on :$(DLV_PORT))
  - $(call command-style,test,     Run all Go tests)
  - $(call command-style,clean,    Delete built artifacts)
  - $(call command-style,[help],   Print this help)
endef
export help_text

.PHONY: test clean help build all install run debug

help:
	@echo -e "$$help_text"

# Select the right binary for the current host
ifeq ($(OS),Windows_NT)
  BIN := $(BIN_DIR)/windows/$(PACKAGE)$(EXTENSION)
else
  UNAME := $(shell uname -s)
  ifeq ($(UNAME),Linux)
    BIN := $(BIN_DIR)/linux/$(PACKAGE)
  endif
  ifeq ($(UNAME),Darwin)
    BIN := $(BIN_DIR)/darwin/$(PACKAGE)
  endif
endif

SOURCES := $(shell find . -name "*.go")
SOURCES += go.mod go.sum

all: $(BINARIES)
	@echo -e "$(GREEN)📦️ Builds are complete: $(END)$(PURPLE)$(BIN_DIR)$(END)"

$(BIN_DIR)/%/$(PACKAGE)$(EXTENSION): $(SOURCES)
	@echo -e "$(YELLOW)🚧 Building $@...$(END)"
	@CGO_ENABLED=1 GOARCH=$(GOARCH) GOOS=$* $(GOBUILD) -o $@ $(BUILD_ENTRY)

build: $(BIN)
	@echo -e "$(GREEN)📦️ Build is complete: $(END)$(PURPLE)$(BIN)$(END)"

clean:
	@rm -rf $(BIN_DIR)
	@echo -e "$(GREEN)Cleaned!$(END)"

TEST_FILES = $(PWD)/internal/
test:
	@echo -e "$(YELLOW)Testing...$(END)"
	@go test $(TEST_FILES)
	@echo -e "$(GREEN)✅ Test is complete!$(END)"

run: $(BIN)
	@exec $?

cpm: $(BIN)
	@exec $? cpm

debug:
	@go run main.go -d

install: $(BIN)
	@echo -e "$(YELLOW)🚀 Installing $(BIN) to $(GOPATH)/bin...$(END)"
	@$(GOINSTALL) $(BUILD_ENTRY)
	@mv $(GOPATH)/bin/$(PACKAGE) $(GOPATH)/bin/$(PACKAGE)
	@echo -e "$(GREEN)✅ Installation complete!$(END)"
