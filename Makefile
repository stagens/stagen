VERSION := 1.0.1
IMAGE := vidog/stagen

GIT_TAG         := $(shell git describe --tags --dirty --always 2>/dev/null || echo "dev")
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT      := $(shell git rev-parse HEAD)
GIT_COMMIT_SHORT:= $(shell git rev-parse --short HEAD)

LDFLAGS := -X 'stagen/internal/build.GitTag=$(GIT_TAG)' \
           -X 'stagen/internal/build.GitBranch=$(GIT_BRANCH)' \
           -X 'stagen/internal/build.GitCommit=$(GIT_COMMIT)' \
           -X 'stagen/internal/build.GitCommitShort=$(GIT_COMMIT_SHORT)'

.PHONY: all
all: dep gen lint build test

.PHONY: cleanup
cleanup:
	rm -rf examples/build

.PHONY: dep
dep:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: gen
gen:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run --tests

.PHONY: go_test
go_test:
	GOMAXPROCS=4 go test ./... -p 4 -parallel 4 -count=1

.PHONY: test
test: go_test

.PHONY: build
build:
	go build \
		-ldflags "$(LDFLAGS)" \
		-o stagen \
		./cmd/stagen/stagen.go

.PHONY: install
install:
	sudo rm -rf /usr/bin/stagen /usr/local/bin/stagen
	sudo ln -s $(CURDIR)/stagen /usr/local/bin/stagen

.PHONY: web
web:
	python3 -m http.server -d example/build 8001

.PHONY: docker_build
docker_build:
	docker buildx build \
		--platform linux/amd64 \
		--build-arg GIT_TAG=$(GIT_TAG) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_COMMIT_SHORT=$(GIT_COMMIT_SHORT) \
		-t $(IMAGE):$(VERSION) \
		--load .

.PHONY: docker_publish
docker_publish:
	docker push $(IMAGE):$(VERSION)

.PHONY: docker_publish_latest
docker_publish_latest:
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):latest

.PHONY: publish
publish: docker_build docker_publish docker_publish_latest
