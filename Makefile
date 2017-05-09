RELEASES=bin/ngx-darwin-amd64 \
	 bin/ngx-linux-amd64 \
	 bin/ngx-linux-386 \
	 bin/ngx-linux-arm \
	 bin/ngx-windows-amd64.exe \
	 bin/ngx-windows-386.exe \
	 bin/ngx-solaris-amd64 

all: $(RELEASES)

bin/ngx-%: GOOS=$(firstword $(subst -, ,$*))
bin/ngx-%: GOARCH=$(subst .exe,,$(word 2,$(subst -, ,$*)))
bin/ngx-%: $(wildcard *.go)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build \
	     -ldflags "-X main.osarch=$(GOOS)/$(GOARCH) -s -w" \
	     -buildmode=exe \
	     -tags release \
	     -o $@

clean:
	rm -rf bin
