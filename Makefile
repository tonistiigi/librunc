export GOPATH:=$(CURDIR)/Godeps/_workspace:$(GOPATH)

lint: vet
	go vet ./...
	go fmt ./...

vet:
	go get golang.org/x/tools/cmd/vet

test:
	go test ${TESTFLAGS} -v ./...


