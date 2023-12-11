MOCKERY_RELEASE=https://github.com/vektra/mockery/releases/download/v2.22.1/mockery_2.22.1_Linux_x86_64.tar.gz

tmp/mockery:
	mkdir -p tmp/
	cd tmp/ && \
	curl -L $(MOCKERY_RELEASE) | tar xvz

mocks: tmp/mockery
	tmp/mockery --all --dir core/

test: unit-test architecture-test

ISETTA_PACKAGES=$(shell go list org.samba/isetta/... | egrep -v 'mocks|org.samba.isetta$$')
unit-test: mocks
	go clean -testcache
	go test -timeout 30s -run ^Test $(ISETTA_PACKAGES) -v -coverprofile tmp/coverage.out

unit-test-with-wsl: mocks
	go clean -testcache
	go test -timeout 60s -tags wsl -run ^Test $(ISETTA_PACKAGES) -v

# additionally requires user to provide admin credentials for elevated execution test
unit-test-with-wsl-interactive: mocks
	go clean -testcache
	go test -timeout 60s -tags wsl,interactive -run ^Test $(ISETTA_PACKAGES) -v

coverage: mocks
	go clean -testcache
	go test -timeout 30s -run ^Test $(ISETTA_PACKAGES) -v -coverprofile=tmp/coverage.out
	go tool cover -func tmp/coverage.out

coverage-with-wsl-interactive: mocks
	go clean -testcache
	go test -timeout 60s -tags wsl,interactive -run ^Test $(ISETTA_PACKAGES) -v -coverprofile=tmp/coverage.out
	go tool cover -func tmp/coverage.out

ARCH_GO=~/go/bin/arch-go
architecture-test: $(ARCH_GO)
	~/go/bin/arch-go -v

$(ARCH_GO): 
	unset GOPATH && \
	go install github.com/fdaines/arch-go

tmp/isetta:
	GOOS=linux go build -ldflags="-w -s" -o tmp/isetta . 

clean:
	rm -rf tmp/ mocks/