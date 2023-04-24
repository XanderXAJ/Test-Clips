.DEFAULT: build
.PHONY: build run

build: build-tc build-tcb

build-tc:
	go build ./cmd/test-clips

build-debug-tc:
	go build -gcflags="all=-N -l" -o test-clips-debug ./cmd/test-clips

debug-tc: build-debug-tc
	dlv exec ./test-clips-debug

build-tcb:
	go build ./cmd/tcb

test:
	go test ./...
