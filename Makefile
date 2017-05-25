########################################################################################

# This Makefile generated by GoMakeGen 0.6.0 using next command:
# gomakegen --metalinter --strip .

########################################################################################

.PHONY = fmt all clean deps metalinter

########################################################################################

all: source-index

source-index:
	go build -ldflags="-s -w" source-index.go

deps:
	git config --global http.https://pkg.re.followRedirects true
	go get -d -v pkg.re/essentialkaos/ek.v9

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

metalinter:
	test -s $(GOPATH)/bin/gometalinter || (go get -u github.com/alecthomas/gometalinter ; $(GOPATH)/bin/gometalinter --install)
	$(GOPATH)/bin/gometalinter --deadline 30s

clean:
	rm -f source-index

########################################################################################
