.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build dev images
	docker-compose build --no-cache

test: ## Run the tests
	docker-compose exec api go test -v -race

run: ## Start the containers and run the spacy worker
	docker-compose up -d --remove-orphans