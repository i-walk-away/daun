.PHONY: build install run test bench clean

BINARY := daun
INSTALL_DIR := $(HOME)/.local/bin

build:
	go build -o $(BINARY) ./cmd/daun

install:
	mkdir -p $(INSTALL_DIR)
	go build -o $(INSTALL_DIR)/$(BINARY) ./cmd/daun

run:
	go run ./cmd/daun

test:
	go test ./...

bench:
	systemd-run --user --scope \
		-p MemoryMax=4G \
		go test ./internal/buffer \
		-run '^$$' \
		-bench . \
		-benchmem \
		-benchtime=2s

clean:
	rm -f $(BINARY)