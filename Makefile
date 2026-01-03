COVERAGE_EXCLUDE=mocks|vendor|test|docs|main.go|config.go|client.go
COVERAGE_THRESHOLD = 80

.PHONY: swag-gen
swag-gen:
	swag init -g ./cmd/api/main.go --output ./docs

.PHONY: run
run: swag-gen
	go run ./cmd/api/main.go

.PHONY: mock-gen
mock-gen:
	go generate ./...

.PHONY: test docker-test
test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -v -E "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi

COVERAGE_FOLDER=./coverage

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark-service-test:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
	else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

.PHONY: redis-run redis-cli redis-monitor
redis-run:
	docker run --name redis -p 6379:6379 -d redis

redis-cli:
	docker exec -it redis redis-cli

redis-monitor:
	docker exec -it redis redis-cli monitor

.PHONY: build

build:
	docker build -t bookmark_service_123:dev .