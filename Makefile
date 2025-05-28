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
	validate-api applications/rest-api/openapi.yaml

.PHONY: rest-api/generate
rest-api/generate: ## (Re)generate generated code for the REST API application
	cd applications/rest-api && go generate

.PHONY: rest-api/database/up
rest-api/database/up: ## Apply the up migrations for the REST API application's database
	psql -h postgres -U postgres -d postgres -f applications/rest-api/internal/database/migrations/up.sql

.PHONY: rest-api/database/down
rest-api/database/down: ## Apply the down migrations for the REST API application's database
	psql -h postgres -U postgres -d postgres -f applications/rest-api/internal/database/migrations/down.sql

.PHONY: rest-api/database/populate
rest-api/database/populate: ## Populate the REST API application's database from a CSV
	psql -h postgres -U postgres -d postgres -c "\copy posts(title,content,created_at) FROM 'applications/rest-api/data.csv' DELIMITER ',' CSV"

.PHONY: rest-api/database/truncate
rest-api/database/truncate: ## Truncates all data in the REST API application's database, but leaves schema unchanged
	psql -h postgres -U postgres -d postgres -f applications/rest-api/truncate.sql

.PHONY: database/connect
database/connect: ## Connect to the postgres database
	psql -h postgres -U postgres -d postgres
