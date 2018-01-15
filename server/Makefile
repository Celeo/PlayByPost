all: run

run: build
	@DATABASE_URL='user=postgres password=postgres host=localhost port=5432 dbname=playbypost sslmode=disable' JOIN_CODE=hello-world ./server

build:
	@go build

clean:
	@go clean

deploy:
	@cd .. && git subtree push --prefix server heroku master
