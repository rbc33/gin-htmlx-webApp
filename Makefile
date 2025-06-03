GOCMD=go
TEMPL=templ
TAILWIND=tailwindcss
BUILD_DIR=./tmp
BINARY_NAME=gocms
ADMIN_BINARY_NAME=gocms-admin

GOCMS_PATH=./cmd/$(BINARY_NAME)
GOCMS_ADMIN_PATH=./cmd/$(ADMIN_BINARY_NAME)

all: build test

build:
	$(TEMPL) generate
	$(TAILWIND) -i ./static/style.css -o ./static/output.css -m 
	$(GOCMD) build -v -o $(BUILD_DIR)/$(BINARY_NAME) $(GOCMS_PATH)
	$(GOCMD) build -v -o $(BUILD_DIR)/$(ADMIN_BINARY_NAME) $(GOCMS_ADMIN_PATH)
test:
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
.PHONY: all build test clean
		