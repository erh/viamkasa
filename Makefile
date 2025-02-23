
GO_BUILD_ENV :=
GO_BUILD_FLAGS :=
MODULE_BINARY := bin/viamkasamodule

ifeq ($(VIAM_TARGET_OS), windows)
	GO_BUILD_ENV += GOOS=windows GOARCH=amd64
	GO_BUILD_FLAGS := -tags no_cgo	
	MODULE_BINARY = bin/viamkasamodule.exe
endif

$(MODULE_BINARY): Makefile go.mod *.go cmd/module/*.go 
	$(GO_BUILD_ENV) go build $(GO_BUILD_FLAGS) -o $(MODULE_BINARY) cmd/module/cmd.go

lint:
	gofmt -s -w .

update:
	go get go.viam.com/rdk@latest
	go mod tidy

test:
	go test ./...

module.tar.gz: meta.json $(MODULE_BINARY) start.sh
	tar czf $@ meta.json $(MODULE_BINARY) start.sh
	git checkout meta.json

ifeq ($(VIAM_TARGET_OS), windows)
module.tar.gz: fix-meta-for-win
else
module.tar.gz: strip-module
endif

strip-module: bin/viamkasamodule
	strip bin/viamkasamodule

fix-meta-for-win:
	jq '.entrypoint = "./bin/viamkasamodule.exe"' meta.json > temp.json && mv temp.json meta.json

module: test module.tar.gz

all: test bin/viamkasa module 

setup:
