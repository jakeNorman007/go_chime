#Makefile
include .env

#run main program
run:
	@echo Starting project...
	@go run main.go

#database commands
up:
	@echo ========== migrating databsse up ==========
	@cd ./db/migrations && goose postgres ${DB_URL} up

down:
	@echo ========== migrating database down ==========
	@cd ./db/migrations && goose postgres ${DB_URL} down
