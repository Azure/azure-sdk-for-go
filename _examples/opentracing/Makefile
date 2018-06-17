# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

M = $(shell printf "\033[34;1m▶\033[0m")
Q = $(if $(filter 1,$V),,@)

all: clean build ; $(info $(M) building opentracing example…) @ ## Build opentracing example
build:
		$Q $(GOBUILD) -o ./bin/opentracing

clean:
		$Q $(GOCLEAN) && rm -rf ./bin