GIT_VER := $(shell git describe --tags)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
export GO111MODULE := on

.PHONY: test binary install clean

cmd/twitter-stream-rule-manager/twitter-stream-rule-manager: *.go cmd/twitter-stream-rule-manager/*.go go.* appspec/*.go
	cd cmd/twitter-stream-rule-manager && go build -ldflags "-s -w -X main.Version=${GIT_VER} -X main.buildDate=${DATE}" -gcflags="-trimpath=${PWD}"

install: cmd/twitter-stream-rule-manager/twitter-stream-rule-manager
	install cmd/twitter-stream-rule-manager/twitter-stream-rule-manager ${GOPATH}/bin

test:
	go test -race ./...

packages:
	cd cmd/twitter-stream-rule-manager && gox -os="linux darwin" -arch="amd64" -output "../../pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-w -s -X main.Version=${GIT_VER} -X main.buildDate=${DATE}"
	cd pkg && find . -name "*${GIT_VER}*" -type f -exec zip {}.zip {} \;

clean:
	rm -f cmd/twitter-stream-rule-manager/twitter-stream-rule-manager
	rm -f pkg/*

release:
	ghr -prerelease -u mashiike -r twitter-stream-rule-manager -n "$(GIT_VER)" $(GIT_VER) pkg/

ci-test:
	$(MAKE) install
	cd tests/ci && PATH=${GOPATH}/bin:$PATH $(MAKE) test
