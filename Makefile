
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

TAG ?= $(USER)
BUILD_LDFLAGS := '$(LDFLAGS)'

all: build


getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golint 1>/dev/null || (echo "Installing golint" && go get -u golang.org/x/lint/golint)
	@which staticcheck 1>/dev/null || (echo "Installing staticcheck" && wget --quiet -O ${GOPATH}/bin/staticcheck https://github.com/dominikh/go-tools/releases/download/2019.1/staticcheck_linux_amd64 && chmod +x ${GOPATH}/bin/staticcheck)
	@which misspell 1>/dev/null || (echo "Installing misspell" && wget --quiet https://github.com/client9/misspell/releases/download/v0.3.4/misspell_0.3.4_linux_64bit.tar.gz && tar xf misspell_0.3.4_linux_64bit.tar.gz && mv misspell ${GOPATH}/bin/misspell && chmod +x ${GOPATH}/bin/misspell && rm -f misspell_0.3.4_linux_64bit.tar.gz)


verifiers: getdeps vet fmt lint staticcheck spelling

vet:
	@echo "Running $@"
	@GO111MODULE=on go vet github.com/zaibon/tfdirectory/...

fmt:
	@echo "Running $@"
	@GO111MODULE=on gofmt -d cmd/
	@GO111MODULE=on gofmt -d mongo/
	@GO111MODULE=on gofmt -d server/
	@GO111MODULE=on gofmt -d *.go

lint:
	@echo "Running $@"
	@GO111MODULE=on ${GOPATH}/bin/golint -set_exit_status github.com/zaibon/tfdirectory/...

staticcheck:
	@echo "Running $@"
	@GO111MODULE=on ${GOPATH}/bin/staticcheck github.com/zaibon/tfdirectory/...

# spelling:
# 	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find cmd/`
# 	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find pkg/`
# 	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find docs/`
# 	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find buildscripts/`
# 	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find dockerscripts/`

# Builds tfdirectory, runs the verifiers then runs the tests.
check: test
test: verifiers build
	@echo "Running unit tests"
	@GO111MODULE=on CGO_ENABLED=0 go test ./... 1>/dev/null

coverage: build
	@echo "Running all coverage for minio"
	@GO111MODULE=on CGO_ENABLED=0 go test -v -coverprofile=coverage.txt -covermode=atomic ./...

# Builds tfdirectory locally.
build:
	@echo "Building tfdirectory binary to './tfdirectory'"
	@GO111MODULE=on GOFLAGS="" CGO_ENABLED=0 go build --ldflags $(BUILD_LDFLAGS) -o $(PWD)/tfdirectory github.com/zaibon/tfdirectory/cmd 1>/dev/null

# docker: build
# 	@docker build -t $(TAG) . -f Dockerfile.dev

# Builds minio and installs it to $GOPATH/bin.
install: build
	@echo "Installing tfdirectory binary to '$(GOPATH)/bin/tfdirectory'"
	@mkdir -p $(GOPATH)/bin && cp -uf $(PWD)/tfdirectory $(GOPATH)/bin/tfdirectory
	@echo "Installation successful. To learn more, try \"tfdirectory --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf tfdirectory
	@rm -rvf build
	@rm -rvf release