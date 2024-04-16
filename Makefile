.PHONY: all
all: help tidy lint test test/cover run

.PHONY: help
help: ## Display this help screen
	awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run: ## run redis_codecrafters
	docker compose up -d --build && docker compose logs -f redis_codecrafters

.PHONY: run_redis
run_redis: ## run redis
	docker exec -it redis redis-cli -h redis_codecrafters

.PHONY: kill
kill: ## kill
	docker rm --force $$(docker ps -aq)