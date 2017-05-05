RELEASES=bin/ngc-darwin-amd64 \
	 bin/ngc-linux-amd64 \
	 bin/ngc-linux-386 \
	 bin/ngc-linux-arm \
	 bin/ngc-windows-amd64.exe \
	 bin/ngc-windows-386.exe \
	 bin/ngc-solaris-amd64 

all: $(RELEASES)

bin/ngc-%: GOOS=$(firstword $(subst -, ,$*))
bin/ngc-%: GOARCH=$(subst .exe,,$(word 2,$(subst -, ,$*)))
bin/ngc-%: $(wildcard *.go)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build \
	     -ldflags "-X main.osarch=$(GOOS)/$(GOARCH) -s -w" \
	     -buildmode=exe \
	     -tags release \
	     -o $@

clean:
	rm -rf bin
