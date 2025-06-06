GOCMD=go
TEMPL=templ
# BUILD_DIR=./
BINARY_NAME=gocms
ADMIN_BINARY_NAME=gocms-admin
GOOSE=goose

GOCMS_PATH=.

DB_DRIVER=mysql
MIGRATIONS_DIR=./migrations

DEV_DB_URL := $(shell grep MY_SQL_URL .env | cut -d'=' -f2-)
PROD_DB_URL := $(shell clever env | grep MYSQL_ADDON_USER | cut -d'=' -f2):$(shell clever env | grep MYSQL_ADDON_PASSWORD | cut -d'=' -f2)@tcp($(shell clever env | grep MYSQL_ADDON_HOST | cut -d'=' -f2):$(shell clever env | grep MYSQL_ADDON_PORT | cut -d'=' -f2))/$(shell clever env | grep MYSQL_ADDON_DB | cut -d'=' -f2)

all: build test

# Development build con CSS
build:
	$(TEMPL) generate
	tailwindcss -i ./static/style.css -o ./static/output.css -m 
	$(GOCMD) build -v -o $(BINARY_NAME) $(GOCMS_PATH)

# Production build (sin tailwind, CSS ya debe existir)
build-prod:
	$(TEMPL) generate
	CGO_ENABLED=0 GOOS=linux $(GOCMD) build -ldflags="-w -s" -o $(BINARY_NAME) $(GOCMS_PATH)

# Generar solo CSS
css:
	tailwindcss -i ./static/style.css -o ./static/output.css -m

# Resto de tus comandos...
test:
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm  ./static/output.css
	rm ./views/tailwind/*.go
	
migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DEV_DB_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DEV_DB_URL)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DEV_DB_URL)" status

migrate-prod:
	@echo "🗄️  Running production migrations on Clever Cloud..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(PROD_DB_URL)" up

migrate-status-prod:
	@echo "📊 Production migration status:"
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(PROD_DB_URL)" status

migration:
	@read -p "Migration name: " name; \
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $$name sql

# Agregar esto a tu Makefile
debug-db:
	@echo "=== Debug Database Connection ==="
	@echo "MYSQL_ADDON_USER: $$(clever env | grep MYSQL_ADDON_USER | cut -d'=' -f2)"
	@echo "MYSQL_ADDON_HOST: $$(clever env | grep MYSQL_ADDON_HOST | cut -d'=' -f2)"
	@echo "MYSQL_ADDON_PORT: $$(clever env | grep MYSQL_ADDON_PORT | cut -d'=' -f2)"
	@echo "MYSQL_ADDON_DB: $$(clever env | grep MYSQL_ADDON_DB | cut -d'=' -f2)"
	@echo "Full URL: $(PROD_DB_URL)"


.PHONY: all build build-prod css test clean migrate-up migrate-down migrate-status migrate-prod migrate-status-prod migration