# import environment variables from ".env" file;
# ".env" file not stored in git, but ".env.example" is stored in git for reference
include ./.env

# makefile is used for automating tasks, not for file dependencies
.PHONY: $(wildcard *)

all: run

# --------------------------------------------------------------------------------

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/  /'
	@echo ""
	@echo "Example development workflow:"
	@echo "   db-setup/db-start ... gen-watch ... [CODING] ... run / db-migrate (*) ... oci-build ... oci-run ... db-stop ... deploy ... [DEPLOYING]" 

# --------------------------------------------------------------------------------

## gen: Generate go code from templ files and tailwindcss style
gen:
	templ generate
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./src/resources/public/style.css --minify

# --------------------------------------------------------------------------------

## gen-watch: Watch and generate go code from templ files and tailwindcss style
gen-watch:
	templ generate --watch &
# alternative for watched generation of css using "npm run dev"
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./src/resources/public/style.css --watch 
# jobs
# kill %1

# --------------------------------------------------------------------------------

## clean: Clean temp/app directory
clean:
	rm -rf temp/app
	mkdir -p ./temp/app

# --------------------------------------------------------------------------------

## build: Build application executable to prepare for oci image creation 
## : (no .env configuration is used)
build: clean gen
	go mod tidy
# GOOS=windows GOARCH=amd64 go build -o temp/$(APP_NAME)-windows.exe ./src
# file temp/$(APP_NAME)-windows.exe
# GOOS=linux GOARCH=amd64 go build -o temp/$(APP_NAME)-linux.exe ./src
# file temp/$(APP_NAME)-linux.exe
	GOOS=linux GOARCH=amd64 go build -o temp/app/main.exe ./src
# Cross-Compilation:
# Install gcc-mingw toolchain on Ubuntu 24 for Windows cross-compilation:
#	sudo apt-get update && sudo apt-get install gcc-mingw-w64
# Install tdm-gcc toolchain on Windows 11 with WSL2 for Linux cross-compilation:
#   # Download and install tdm-gcc from https://jmeubank.github.io/tdm-gcc/

# --------------------------------------------------------------------------------

## run (default): Run application from source code for development purposes
## : (database must be accessible and .env configuration is used)
run: 
# we assume that tailwindcss style and go code from templ files are already generated
	go run ./src

# --------------------------------------------------------------------------------

# current time for labeling docker images with build time;
MAKE_TIME = $(shell date +"%FT%H:%M:%SZ")

## oci-build: Build application oci image
oci-build: build
	docker build . -t $(APP_NAME):latest --label "version=${APP_VERSION}" --label "build=$(MAKE_TIME)" \
	  --label "maintainer=https://github.com/RobertTC32/app-$(APP_NAME)"
	docker image ls | grep $(APP_NAME) 
	docker image inspect $(APP_NAME):latest
# docker rmi $(APP_NAME):latest

# --------------------------------------------------------------------------------

## oci-run: Run existing application oci image for development purposes
## : (database must be accessible and .env configuration is used)
oci-run: 
	docker run --name $(APP_NAME) --rm -p $(APP_PORT):80 --env-file ./.env $(APP_NAME):latest

# --------------------------------------------------------------------------------

## db-clean: Clean temp/db directory
db-clean:
	rm -rf temp/db
	mkdir -p ./temp/db

# --------------------------------------------------------------------------------

## db-start: Start existing development database
## : (.env configuration is used)
db-start:  
	# start mariadb container with mounted volumes for database and migration files
	docker run --name db-$(APP_NAME) --rm -d -v ./temp/db:/var/lib/mysql:Z -p $(OCI_DB_PORT):3306 --env-file ./.env mariadb:12.2 
# docker exec -it db-$(APP_NAME) mariadb -uroot -p

# --------------------------------------------------------------------------------

## db-migrate: Migrate existing development database
## : (database must be accessible and .env configuration is used)
db-migrate:
# migrate database using goose sql scripts in src-db/migrations/ directory
# (create goose sql scripts using "goose -s create ${name} sql" command)
	goose up 

# --------------------------------------------------------------------------------

## db-setup: Create, start and migrate new development database
## : (.env configuration is used)
db-setup: db-clean db-start db-migrate  
# create development database in temp/db directory from scratch
# Start mariadb container with mounted volumes for database 
# Execute migration files on database

# --------------------------------------------------------------------------------

## db-stop: Stop running development database
## : (database must be running and no .env configuration is used)
db-stop:
	docker stop db-$(APP_NAME) 

# --------------------------------------------------------------------------------

## deploy: Create deployment package with application oci image included
deploy: oci-build
	rm -rf temp/deploy-$(APP_NAME)
	mkdir -p ./temp/deploy-$(APP_NAME)
	# copy all from src-deploy
	cp -r ./src-deploy/scripts/. ./temp/deploy-$(APP_NAME)/
	# copy all from src-db
	mkdir -p ./temp/deploy-$(APP_NAME)/migrations
	cp -r ./src-db/migrations/. ./temp/deploy-$(APP_NAME)/migrations/
	# copy image from local oci registry
	docker save -o ./temp/deploy-$(APP_NAME)/$(APP_NAME).tar $(APP_NAME):latest
	du -sh ./temp/deploy-$(APP_NAME)/$(APP_NAME).tar
# docker load -i ./temp/deploy-$(APP_NAME)/$(APP_NAME).tar
	# copy goose from containing folder
	cp ../goose.exe ./temp/deploy-$(APP_NAME)/

# --------------------------------------------------------------------------------
