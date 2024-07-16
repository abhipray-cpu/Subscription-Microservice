SHELL=/bin/bash
SUBSCRIPTION_BINARY=subscriptionApp
LISTENER_BINARY=listenerApp
LOGGER_BINARY=loggerApp
PAYMENT_BINARY=paymentApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_subscription build_listener build_logger build_payment
	@echo "Stopping docker images (if running)..."
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:	
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_subscription: builds the subscription binary as a linux executable
build_subscription:
	@echo "Building subscription binary..."
	cd ../subscription-service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o ${SUBSCRIPTION_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o ${LISTENER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_payment: builds the payment binary as a linux executable
build_payment:
	@echo "Building payment binary..."
	cd ../payment-service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o ${PAYMENT_BINARY} ./cmd/api
	@echo "Done!"