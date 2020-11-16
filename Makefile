.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: dep
dep: ## Установка и фиксация зависимостей
	$(eval PACKAGE := $(shell go list -m))
	@go mod download
	@go mod vendor

.PHONY: test
test: dep ## Запуск юнит тестов
	@go test ./...

.PHONY: build
build: test ## Сборка бинарников
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/distate-task ./cmd/distate-task

.PHONY: docker
docker: ## Создание docker образа
	docker build -t distate-task -f ./Dockerfile .

.PHONY: up
up: build docker ## Поднятие контейнеров
	docker-compose up -d

.PHONY: migrate
migrate: ## Запуск миграций
	tern migrate --config ./config/tern.conf --migrations ./migrations

.PHONY: down
down: ## Остановка контейнеров
	docker-compose down

.PHONY: clean
clean: ## Очистка от контейнеров и образов
	-docker container stop $$(docker ps -q -a)
	-docker container rm $$(docker ps -q -a)
	-docker image rm -f $$(docker image ls -q)
