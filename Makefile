.DEFAULT: build
.PHONY: build run

build:
	go build .

test:
	go test .
