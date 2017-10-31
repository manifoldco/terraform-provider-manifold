VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
PKG=github.com/manifoldco/terraform-provider-manifold
GO_BUILD=CGO_ENABLED=0 go build -i --ldflags="-w -X $(PKG)/config.Version=$(VERSION)"
PROMULGATE_VERSION=0.0.8

LINTERS=\
    gofmt \
    golint \
    gosimple \
    vet \
    misspell \
    ineffassign \
    deadcode

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
    github.com/golang/dep/cmd/dep \
    github.com/alecthomas/gometalinter \
    github.com/jbowes/oag

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	gometalinter --install

vendor: Gopkg.lock
	dep ensure

.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Test and linting
#################################################

test: vendor
	@CGO_ENABLED=0 go test -v ./...

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
     ./... --enable

$(LINTERS): vendor
	$(METALINT) $@

.PHONY: $(LINTERS) test

#################################################
# Building
#################################################

PREFIX?=
SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

build: $(PREFIX)bin/terraform-provider-manifold$(SUFFIX)

$(PREFIX)bin/terraform-provider-manifold$(SUFFIX):
	$(GO_BUILD) -o $(PREFIX)bin/terraform-provider-manifold$(SUFFIX) .

.PHONY: build

NO_WINDOWS= \
	darwin_amd64 \
	linux_amd64
OS_ARCH= \
	$(NO_WINDOWS) \
	windows_amd64

os=$(word 1,$(subst _, ,$1))
arch=$(word 2,$(subst _, ,$1))

os-build/windows_amd64/bin/terraform-provider-manifold: os-build/%/bin/terraform-provider-manifold:
	PREFIX=build/$*/ GOOS=$(call os,$*) GOARCH=$(call arch,$*) make build/$*/bin/terraform-provider-manifold.exe
$(NO_WINDOWS:%=os-build/%/bin/terraform-provider-manifold): os-build/%/bin/terraform-provider-manifold:
	PREFIX=build/$*/ GOOS=$(call os,$*) GOARCH=$(call arch,$*) make build/$*/bin/terraform-provider-manifold

build/terraform-provider-manifold_$(VERSION)_windows_amd64.zip: build/terraform-provider-manifold_$(VERSION)_%.zip: os-build/%/bin/terraform-provider-manifold
	cd build/$*/bin; zip -r ../../terraform-provider-manifold_$(VERSION)_$*.zip terraform-provider-manifold.exe
$(NO_WINDOWS:%=build/terraform-provider-manifold_$(VERSION)_%.tar.gz): build/terraform-provider-manifold_$(VERSION)_%.tar.gz: os-build/%/bin/terraform-provider-manifold
	cd build/$*/bin; tar -czf ../../terraform-provider-manifold_$(VERSION)_$*.tar.gz terraform-provider-manifold

zips: $(NO_WINDOWS:%=build/terraform-provider-manifold_$(VERSION)_%.tar.gz) build/terraform-provider-manifold_$(VERSION)_windows_amd64.zip

release: zips
	curl -LO https://releases.manifold.co/promulgate/$(PROMULGATE_VERSION)/promulgate_$(PROMULGATE_VERSION)_linux_amd64.tar.gz
	tar xvf promulgate_*
	./promulgate release v$(VERSION)

.PHONY: release zips $(OS_ARCH:%=os-build/%/bin/manifold)
