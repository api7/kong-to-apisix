test:
	./tools/setup.sh
	go test -v -count=1 -race -cover -coverprofile=coverage.txt -covermode=atomic ./...