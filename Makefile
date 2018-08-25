all: build

test: build
	go $@ -v -cover -coverprofile=count.out

build: deps
	go build

.PHONY: clean
clean:
	go clean && rm -f count.out

deps:
	go get -d ./...

coverage: test
	sed -i "s/.*\//.\//" count.out && go tool cover -html=count.out
