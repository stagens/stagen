# Yet Another Static Website Generator

## Installation

### Requirements
- go 1.25+
- git

### Optional requirements
- [pagefind](https://pagefind.app/) (for website search without backend)
- docker (for using pre-built docker image)

```shell
$ go install github.com/stagens/stagen/cmd/stagen@latest
```

## Build from sources

```shell
$ git clone https://github.com/stagens/stagen.git
$ cd stagen
$ make
$ make install
```

## Examples

### Create new project in current directory

```shell
$ stagen init .
```

### Build

```shell
$ stagen build
``` 

### Dev (start dev server and rebuild changes)

```shell
$ stagen dev
```

## Using docker

### Download docker image

```shell
$ docker pull vidog/stagen
```

### Create new project

```shell
$ docker run -ti --rm -v'./:/project' vidog/stagen init /project
```

### Add default theme

```shell
$ git init
$ git submodule add https://github.com/stagens/theme-default.git themes/default
```

### Build project

```shell
$ docker run -ti --rm -v'./:/project' vidog/stagen build /project
```

### Project development

```shell
$ docker run -ti --rm -v'./:/project' -p8001:8001 vidog/stagen dev /project
```

## Docs

TBD
