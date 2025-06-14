GOCMD=go
TEMPL=templ
TAILWIND=tailwindcss
BUILD_DIR=.
BINARY_NAME=gocms
ADMIN_BINARY_NAME=gocms-admin

GOCMS_PATH=./$(BINARY_NAME)
GOCMS_ADMIN_PATH=./$(ADMIN_BINARY_NAME)

all: build test

prepare_env:
	cp -r migrations tests/system_tests/helpers/

build: prepare_env
	$(TEMPL) generate
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

install-tailwindcss:
	if [ ! -f tailwindcss ]; then \
		wget -q https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.16/tailwindcss-linux-x64 \
			&& echo "33f254b54c8754f16efbe2be1de38ca25192630dc36f164595a770d4bbf4d893  tailwindcss-linux-x64" | sha256sum -c \
			&& chmod +x tailwindcss-linux-x64 \
			&& mv tailwindcss-linux-x64 tailwindcss; \
	fi


.PHONY: all build test clean
		