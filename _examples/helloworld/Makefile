# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

M = $(shell printf "\033[34;1m▶\033[0m")
Q = $(if $(filter 1,$V),,@)

all: clean build ; $(info $(M) building helloworld example…) @ ## Build helloworld example
build:
		$Q $(GOBUILD) -o ./bin/consumer ./consumer/main.go && $(GOBUILD) -o ./bin/producer ./producer/main.go
clean:
		$Q $(GOCLEAN) && rm -rf ./bin