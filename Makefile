.PHONY: test
test:
	ginkgo -r --progress --trace -cover -coverprofile=coverage.txt -covermode=atomic

.PHONY: build
build:
	go build -o bin/kta ./cmd/kong-to-apisix/main.go
