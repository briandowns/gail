GO ?= go
GLIDE ?= glide

BINDIR := bin
BINARY := gail

VERSION := 0.1
LDFLAGS = -ldflags "-X main.gitSHA=$(shell git rev-parse HEAD) -X main.version=$(VERSION) -X main.name=$(BINARY)"

.PHONY:
build: clean
	if [ ! -d $(BINDIR) ]; then mkdir $(BINDIR); fi
	$(GO) build -v -o $(BINDIR)/$(BINARY) $(LDFLAGS)

.PHONY:
test:
	$(GO) test -v -cover .

.PHONY:
clean:
	$(GO) clean
	rm -f $(BINDIR)/$(BINARY)
