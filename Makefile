.PHONY: help
.DEFAULT_GOAL := help

-include .env

swagger: ## generate swagger doc
	swag init -g cmd/main.go --output docs --parseDependency --parseInternal --parseDepth 2
