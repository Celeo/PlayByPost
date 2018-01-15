all: run

run: build
	@./server

build:
	@go build

deploy:
	@cd .. && git subtree push --prefix server heroku master
