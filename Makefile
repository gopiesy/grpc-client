GO_PKG_DIRS  := $(subst $(shell go list -e -m),.,$(shell go list ./ | grep -v /vendor/ ))

all: clean fmt lint
	go build -ldflags="-s -w" -o client $(GO_PKG_DIRS)

fmt:
	gofmt -s -w $(GO_PKG_DIRS)

lint:
	golangci-lint run -v $(GO_PKG_DIRS)

clean:
	rm -f client