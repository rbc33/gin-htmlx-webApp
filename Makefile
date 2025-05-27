GOCMD=go
TEMPL=templ
BUILD_DIR=./tmp
BINARY_NAME=gocms
GOCMS_PATH=./cmd/$(BINARY_NAME)

all: build test

build:
	$(TEMPL) generate
	$(GOCMD) build -v -o $(BUILD_DIR)/$(BINARY_NAME) $(GOCMS_PATH)
test:
	$(GOCMD) test -v ./...
run:
	$(GOCMD) run -v -o $(BUILD_DIR)/$(BINARY_NAME)
clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
		