# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: berith all test clean
.PHONY: berith-linux berith-linux-386 berith-linux-amd64 berith-linux-mips64 berith-linux-mips64le
.PHONY: berith-linux-arm berith-linux-arm-5 berith-linux-arm-6 berith-linux-arm-7 berith-linux-arm64
.PHONY: berith-darwin berith-darwin-386 berith-darwin-amd64
.PHONY: berith-windows berith-windows-386 berith-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

berith:
	build/env.sh go run build/ci.go install ./cmd/berith
	@echo "Done building."
	@echo "Run \"$(GOBIN)/berith\" to launch berith."

# will added after catch compile errors
all:
	build/env.sh go run build/ci.go install

test:
	build/env.sh go run build/ci.go test

test-seq:
	build/env.sh go run build/ci.go test -p 1

test-datasync:
	build/env.sh go run build/ci.go test -p 1 ./datasync/...

test-networks:
	build/env.sh go run build/ci.go test -p 1 ./networks/...

test-tests:
	build/env.sh go run build/ci.go test -p 1 ./tests/...

test-others:
	build/env.sh go run build/ci.go test -p 1 -exclude datasync,networks,tests

fmt:
	build/env.sh go run build/ci.go fmt

lint:
	build/env.sh go run build/ci.go lint

lint-try:
	build/env.sh go run build/ci.go lint-try

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/* build/_workspace/src/

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)
berith-cross: berith-linux berith-darwin berith-windows
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/berith-*

berith-linux: berith-linux-386 berith-linux-amd64 berith-linux-arm berith-linux-mips64 berith-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-*

berith-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/berith
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep 386

berith-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/berith
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep amd64

berith-linux-arm: berith-linux-arm-5 berith-linux-arm-6 berith-linux-arm-7 berith-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep arm

berith-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/berith
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep arm-5

berith-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/berith
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep arm-6

berith-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/berith
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep arm-7

berith-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/berith
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep arm64

berith-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/berith
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep mips

berith-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/berith
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep mipsle

berith-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/berith
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep mips64

berith-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/berith
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/berith-linux-* | grep mips64le

berith-darwin: berith-darwin-386 berith-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/berith-darwin-*

berith-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/berith
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/berith-darwin-* | grep 386

berith-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/berith
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/berith-darwin-* | grep amd64

berith-windows: berith-windows-386 berith-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/berith-windows-*

berith-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/berith
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/berith-windows-* | grep 386

berith-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/berith
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/berith-windows-* | grep amd64
