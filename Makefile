.PHONY: install run test clean

BINARY := daun
INSTALL_DIR := $(HOME)/.local/bin

install:
	mkdir -p $(INSTALL_DIR)
	go build -o $(INSTALL_DIR)/$(BINARY) ./cmd/daun

run:
	go run ./cmd/daun

test:
	go test ./...

clean:
	rm -f $(INSTALL_DIR)/$(BINARY)