# Docdocgo
[![Go Reference](https://pkg.go.dev/badge/github.com/calimapp/docdocgo.svg)](https://pkg.go.dev/github.com/calimapp/docdocgo)

## Install

### With go install
```sh
go install github.com/calimapp/docdocgo@latest
```

### From released binaries

<https://github.com/calimapp/docdocgo/releases>

### With docker

```sh
docker run -v $(pwd):/godoc -w /godoc ghcr.io/calimapp/docdocgo:latest --help
```

## Features

### High priority
- [ ] Dockerize CLI, with auto docker deployment to docker hub
- [ ] Refactor app in 2 separate modules parser and render

### Low priority

- [ ] switch display internal packages and private functions, vars, ...
- [ ] Generate documentation in other format (pdf, markdown, asciidoc, ...)
- [ ] Templates management for html (maybe for other format)
- [ ] Split package doc in different html files ?
