# Lightweight Auto Restart Tool for Golang

Utility `gouto` (golang auto) - very simple and lightweight auto restart tool

## Installation

To install via `go get` (needs `golang` version 1.11+ installed):

```
go get github.com/enfipy/gouto
```

To verify that `gouto` was installed correctly:

```
gouto -h
```

## Usage:

To start auto restart on change simply run:

```
gouto -dir=sources/
```

## Options:

The following is given by running `gouto -h`:

```
Usage of gouto:
  -build string
        Command to rebuild after changes (default "go build")
  -dir string
        Directory to watch for changes (default "./")
  -out string
        Output directory for binary after build (default "./cmd/app")
```
