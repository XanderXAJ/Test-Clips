.DEFAULT: build
.PHONY: build run

build:
	go build ./cmd/test-clips

build-debug:
	go build -gcflags="all=-N -l" -o test-clips-debug ./cmd/test-clips

debug: build-debug
	dlv exec ./test-clips-debug

test:
	go test ./...
