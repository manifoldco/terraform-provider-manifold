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

install: vendor
	go install
	cp $(GOPATH)/bin/terraform-provider-manifold $(HOME)/.terraform.d/plugins/terraform-provider-manifold

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
