GOCMD=go
TEMPL=templ
TAILWIND=tailwindcss
BUILD_DIR=./tmp
BINARY_NAME=gocms
ADMIN_BINARY_NAME=gocms-admin

GOCMS_PATH=./cmd/$(BINARY_NAME)
GOCMS_ADMIN_PATH=./cmd/$(ADMIN_BINARY_NAME)

all: build test

prepare_env:
	cp -r migrations tests/system_tests/helpers/

build: prepare_env
	$(TEMPL) generate
	$(TAILWIND) -i ./static/style.css -o ./static/output.css -m 
	$(GOCMD) build -v -o $(BUILD_DIR)/$(BINARY_NAME) $(GOCMS_PATH)
	$(GOCMD) build -v -o $(BUILD_DIR)/$(ADMIN_BINARY_NAME) $(GOCMS_ADMIN_PATH)

test: prepare_env
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)

install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
	go install github.com/a-h/templ/cmd/templ@v0.3.865
	go install github.com/cosmtrek/air@v1.61.7 

.PHONY: all build test clean
		