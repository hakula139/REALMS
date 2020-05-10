CLI_SRC := ./cmd/realms
DAE_SRC := ./cmd/realmsd
SRCS := $(CLI_SRC) $(DAE_SRC)
BIN_PATH := ./bin

.PHONY: build

build:
	@go build -o $(BIN_PATH) $(CLI_SRC)
	@echo build: realms done.
	@go build -o $(BIN_PATH) $(DAE_SRC)
	@echo build: realmsd done.
