export GOPATH:=$(CURDIR)/Godeps/_workspace:$(GOPATH)

lint: vet
	go vet ./...
	go fmt ./...

vet:
	go get golang.org/x/tools/cmd/vet

test: fixtures/busybox
	go test ${TESTFLAGS} -v ./...

fixtures/busybox:
	mkdir -p $@
	curl -sSL 'https://github.com/jpetazzo/docker-busybox/raw/buildroot-2014.11/rootfs.tar' | tar -xC $@

