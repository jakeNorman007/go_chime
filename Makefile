#Makefile
include .env

#run main program
run:
	@go run main.go

#database commands
up:
	@echo migrating up
	cd ./db/migrations && goose postgres ${DB_URL} up

down:
	@echo migrating down
	cd ./db/migrations && goose postgres ${DB_URL} down
