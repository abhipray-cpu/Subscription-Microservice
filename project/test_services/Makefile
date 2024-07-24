SHELL=/bin/bash

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"


## down: stop current containes and  build new containers
up_build:	
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"
	@echo "Starting docker compose	..."
	docker-compose up --build -d
	@echo "Docker images started!"

## down: stop docker compose
down:	
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"