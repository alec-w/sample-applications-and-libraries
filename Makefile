DEV_ENV_DOCKER_COMPOSE_YAML := dev-env/docker-compose.yaml

# Adaption of https://gist.github.com/prwhite/8168133?permalink_comment_id=4376235#gistcomment-4376235 to allow forward slashes in targets too
.PHONY: help
help: ## Display this help text
	@sed \
		-e '/^[a-zA-Z0-9_/\-]*:.*##/!d' \
		-e 's/:.*##\s*/:/' \
		-e 's/^\(.\+\):\(.*\)/$(shell tput setaf 6)\1$(shell tput sgr0):\2/' \
		$(MAKEFILE_LIST) | column -c2 -t -s :

.PHONY: dev-env/build
dev-env/build: ## Build development environment container image
	docker compose -f $(DEV_ENV_DOCKER_COMPOSE_YAML) build

.PHONY: dev-env/up
dev-env/up: ## Bring the development environment up
	docker compose -f $(DEV_ENV_DOCKER_COMPOSE_YAML) up -d

.PHONY: dev-env/down
dev-env/down: ## Tear the development environment down
	docker compose -f $(DEV_ENV_DOCKER_COMPOSE_YAML) down -v

.PHONY: dev-env/shell
dev-env/shell: ## Start a shell in the development environment
	docker compose -f $(DEV_ENV_DOCKER_COMPOSE_YAML) exec dev bash

.PHONY: rest-api/validate-spec
rest-api/validate-spec: ## Validates REST API application's OpenAPI spec
	validate-api applications/rest-api/schema.yaml

.PHONY: rest-api/generate
rest-api/generate: ## (Re)generate generated code for the REST API application
	cd applications/rest-api && go generate
