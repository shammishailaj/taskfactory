build:
	-docker rm taskfactory
	-docker rmi -f taskfactory:latest
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) docker compose -f ./deployments/docker-compose.yml build --no-cache taskfactory

run:
	# https://stackoverflow.com/a/2670143/6670698
	-docker rm taskfactory_dev
	-docker rmi -f taskfactory_dev:latest
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) docker compose -f ./deployments/docker-compose.yml up --remove-orphans taskfactory-dev

cli:
	sh ./deployments/build_cli.sh $(os) $(arch)

test:
	docker compose -f ./deployments/docker-compose.yml up tests

lint:
	docker compose -f ./deployments/docker-compose.yml up linter

down:
	docker compose -f ./deployments/docker-compose.yml down