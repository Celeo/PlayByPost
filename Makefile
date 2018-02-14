all: deploy

deploy:
	@echo ">>> Building ..."
	@cd server && make build
	@cd client && yarn build

	@echo ">>> Stopping app ..."
	@ssh playbypost "sudo systemctl stop pbp-app"

	@echo ">>> Deploying ..."
	@scp server/server playbypost:/srv/server
	@rsync -avz client/dist/* playbypost:/srv/client/

	@echo ">>> Starting app ..."
	@ssh playbypost "sudo systemctl start pbp-app"

	@echo ">>> Done"
