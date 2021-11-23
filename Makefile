.PHONY: test* run build

PACKAGE_NAME := todo

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deploy-dev: ## deploy to development
	kubectl apply -f infra/dev

deploy-prod: ## deploy to production
	kubectl apply -f infra/prod

build: ## build service
	go build -ldflags="-w -s" -o ./build/app

linuxBuild: ## build service for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/app

run: ## Run service
	go run main.go

test: ## Run all tests
	go test -v ./...

test-race: ## Run all test with race detector
	go test -race -v ./...

clean-test-cache: ## Clean test cache result
	go clean -testcache

coverage: ## Run tests and generate coverage files per package
	mkdir .coverage 2> /dev/null || true
	rm -rf .coverage/*.out || true
	go test ./... -coverprofile=coverage.out -covermode=atomic

clean: ## Clean coverage and build service binary
	rm -rf .coverage/ build/

jaeger: ## Runs Jaeger for observability
	docker run \
	--rm \
	-p5775:5775/udp \
	-p6831:6831/udp \
	-p6832:6832/udp \
	-p5778:5778 \
	-p16686:16686 \
	-p14268:14268 \
	-p14250:14250 \
	-p9411:9411 \
    jaegertracing/all-in-one:1.28