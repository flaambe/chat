GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOBASE := $(shell pwd)
GOBIN=$(GOBASE)/bin
BINARY_NAME=chatserver
    
all: test build
build: 
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v ./cmd
test:
	$(GOTEST) ./...
clean:
	$(GOCLEAN)
	rm -rf $(GOBIN)
run:
	$(GOBIN)/$(BINARY_NAME)
    