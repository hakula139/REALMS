SRC_PATH := cmd/realmsd/main.go
ifeq ($(OS), Windows_NT)
	BINARY_PATH := bin\realmsd.exe
else
	BINARY_PATH := bin/realmsd
endif

.PHONY: build run clean

build:
	@go build -o $(BINARY_PATH) $(SRC_PATH)
	@echo build: done.

run:
	@./$(BINARY_PATH)

clean:
	@go clean
ifeq ($(OS), Windows_NT)
	@del /Q /F $(BINARY_PATH)
else
	@rm -f $(BINARY_PATH)
endif
	@echo clean: done.
