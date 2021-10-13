.PHONY: test
test:
	ginkgo -r --progress --trace -coverpkg=./... -coverprofile=coverage.txt

.PHONY: build
build:
	go build -o bin/kong-to-apisix ./cmd/kong-to-apisix/main.go

.PHONY: unit-test
unit-test:
	go test -race --count=1 ./pkg/...

.PHONY: lint
lint:
	golangci-lint run --verbose ./...
