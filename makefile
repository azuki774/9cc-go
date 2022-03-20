.PHONY: build
build:
	go build -o ./build/9cc-go .

.PHONY: test
test:
	gofmt -l -w .
	go test ./... -v -cover
	test/run.sh

