########################################################################################

.PHONY = fmt all clean deps

########################################################################################

all: source-index

source-index:
	go build -ldflags="-s -w" source-index.go

deps:
	go get -v pkg.re/essentialkaos/ek.v7

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

metalinter:
	test -s $(GOPATH)/bin/gometalinter || go get -u github.com/alecthomas/gometalinter ; $(GOPATH)/bin/gometalinter --install
	$(GOPATH)/bin/gometalinter --deadline 30s

clean:
	rm -f source-index

########################################################################################

