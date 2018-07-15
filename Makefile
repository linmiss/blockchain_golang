BINARY := blockchain

build:
	@echo "====> Go build"
	@go build -o ./build/$(BINARY) ./src

run:
	@echo "-----> Go run"
	./$(BINARY)

.PHONY: build run