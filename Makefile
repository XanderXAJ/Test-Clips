.DEFAULT: build
.PHONY: build run

build:
	go build .

build-debug:
	go build -gcflags="all=-N -l" -o test-clips-debug .

debug: build-debug
	dlv exec ./test-clips-debug

test:
	go test .
