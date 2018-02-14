all: deploy

deploy: build stop_app do_deploy start_app
	@echo ">>> Done"

build:
	@echo ">>> Building ..."
	@cd server && make build
	@cd client && yarn build

stop_app:
	@echo ">>> Stopping app ..."
	@ssh playbypost "sudo systemctl stop pbp-app"

do_deploy:
	@echo ">>> Deploying ..."
	@scp server/server playbypost:/srv/server
	@rsync -avz client/dist/* playbypost:/srv/client/

start_app:
	@echo ">>> Starting app ..."
	@ssh playbypost "sudo systemctl start pbp-app"
