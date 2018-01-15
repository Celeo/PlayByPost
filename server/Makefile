all: run

run: build
	@./server

build:
	@go build

clean:
	@go clean

deploy:
	@cd .. && git subtree push --prefix server heroku master
