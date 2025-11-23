.PHONY: all
all: dep gen lint test build

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

.PHONY: test
test:
	GOMAXPROCS=4 go test ./... -p 4 -parallel 4 -count=1

.PHONY: build
build:
	go build -o stagen ./cmd/stagen/stagen.go

.PHONY: web
web:
	python3 -m http.server -d example/build 8001
