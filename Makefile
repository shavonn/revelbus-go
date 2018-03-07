GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=revelbus
GOOS=linux
GOARCH=amd64
    
all: test build
build: 
		GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o $(BINARY_NAME) -v cmd/main.go
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v cmd/main.go
		./$(BINARY_NAME)