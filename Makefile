DAE_SRC := ./cmd/realmsd
SRCS := $(DAE_SRC)
BIN_PATH := ./bin

.PHONY: build

build:
	@go build -o $(BIN_PATH) $(DAE_SRC)
	@echo build: realmsd done.
