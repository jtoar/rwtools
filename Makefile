clean:
	rm rwtools

build:
	go build

smoke-test: build
	rwtools smoke-test

install:
	go install