.PHONY: test
test:
	ginkgo -r --v --progress --trace -cover -coverprofile=coverage.txt -covermode=atomic

.PHONY: build
build:
	go build -o bin/kta ./cmd/kong-to-apisix/main.go
