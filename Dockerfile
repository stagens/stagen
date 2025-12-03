FROM golang:1.25 AS build

ARG GIT_TAG
ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GIT_COMMIT_SHORT

WORKDIR /app

COPY . .

RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -a \
      -installsuffix cgo \
      -ldflags "\
        -extldflags '-static' \
        -X 'stagen/internal/build.GitTag=${GIT_TAG}' \
        -X 'stagen/internal/build.GitBranch=${GIT_BRANCH}' \
        -X 'stagen/internal/build.GitCommit=${GIT_COMMIT}' \
        -X 'stagen/internal/build.GitCommitShort=${GIT_COMMIT_SHORT}'" \
      -o stagen \
      ./cmd/stagen/stagen.go

FROM alpine AS pagefind

WORKDIR /app

RUN apk add --no-cache wget
RUN apk add --no-cache tar
RUN apk add --no-cache git
RUN wget https://github.com/CloudCannon/pagefind/releases/download/v1.3.0/pagefind-v1.3.0-x86_64-unknown-linux-musl.tar.gz
RUN tar -xf pagefind-v1.3.0-x86_64-unknown-linux-musl.tar.gz

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/stagen /stagen
COPY --from=pagefind /app/pagefind /pagefind

WORKDIR /app

ENTRYPOINT ["/stagen"]
