build:
	docker build -t todo-service .

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

.PHONY: build up down restart