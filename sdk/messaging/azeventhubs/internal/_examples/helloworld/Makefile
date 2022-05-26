# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

all: clean build
build:
		$(GOBUILD) -o ./bin/consumer ./consumer/main.go
		$(GOBUILD) -o ./bin/producer ./producer/main.go
clean:
		$(GOCLEAN)
		rm -rf ./bin