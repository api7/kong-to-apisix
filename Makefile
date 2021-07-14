.PHONY: test
test:
	ginkgo -r --progress --trace -coverpkg=./... -coverprofile=coverage.txt

.PHONY: build
build:
	go build -o bin/kta ./cmd/kong-to-apisix/main.go
