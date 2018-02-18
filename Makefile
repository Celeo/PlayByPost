all: deploy

help:
	@echo "Targets:"
	@echo "(default) = deploy"
	@echo "deploy: builds and deploys both client and server"
	@echo "build_both: builds both client and server"
	@echo "build_client: builds the client"
	@echo "build_server: builds the server"
	@echo "deploy_both: deploys both client and server"
	@echo "do_deploy_client: deploys the client"
	@echo "do_deploy_server: stops the server app, deploys the server, and starts the server app agaion"
	@echo "do_deploy_server_actual: deploys the server app (requires the app to be stopped on the server)"
	@echo "stop_app: stops the server app running on the remote machine"
	@echo "start_app: starts the server app running on the remote machine"

deploy: build_both deploy_both
	@echo ">>> Done"

build_both: build_client build_server

build_client:
	@echo ">>> Building client"
	@cd client && yarn build

build_server:
	@echo ">>> Building server"
	@cd server && make build

deploy_both: do_deploy_client do_deploy_server

do_deploy_client:
	@echo ">>> Deploying client"
	@rsync -avz client/dist/* playbypost:/srv/client/

do_deploy_server: stop_app do_deploy_server_actual start_app

do_deploy_server_actual:
	@echo ">>> Deploying server"
	@scp server/server playbypost:/srv/server

stop_app:
	@echo ">>> Stopping app ..."
	@ssh playbypost "sudo systemctl stop pbp-app"

start_app:
	@echo ">>> Starting app ..."
	@ssh playbypost "sudo systemctl start pbp-app"
