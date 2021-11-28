# dirmap

:file_folder: `dirmap` is a tool for generating a directory map. 

It extracts a part of the document from markdown or source code of each directory and uses it as overview of the directory.

## Quick Start

``` console
$ pwd
/Users/k1low/src/github.com/k1LoW/dirmap
$ dirmap generate
.
├── .github/
│   └── workflows/
├── cmd/ ... Commands.
├── config/ ... Configuration file.
├── matcher/ ... Implementation to find the string that will be the overview from the code or Markdown.
├── output/ ... Output format of the directory map.
├── scanner/ ... Implementation of scanning the target directory and its overview from the file system based on the configuration.
├── scripts/ ... scripts for Dockerfile.
└── version/ ... Version.

$ dirmap generate -t table
| Directory | Overview |
| --- | --- |
| .github/ |  |
| .github/workflows/ |  |
| cmd/ | Commands. ( [ref](cmd/doc.go) ) |
| config/ | Configuration file ( [ref](config/config.go) ) |
| matcher/ | Implementation to find the string that will be the overview from the code or Markdown. ( [ref](matcher/matcher.go) ) |
| output/ | Output format of the directory map ( [ref](output/output.go) ) |
| scanner/ | Implementation of scanning the target directory and its overview from the file system based on the configuration. ( [ref](scanner/scanner.go) ) |
| scripts/ | Scripts for Dockerfile ( [ref](scripts/.dirmap.md) ) |
| version/ | Version ( [ref](version/version.go) ) |
```

## Usage

### `dirmap generate`

`dirmap` collects overviews for each directory and formats them for display.

``` console
$ dirmap generate
.
├── .github/
│   └── workflows/
├── cmd/ ... Commands.
├── config/ ... Configuration file.
├── matcher/ ... Implementation to find the string that will be the overview from the code or Markdown.
├── output/ ... Output format of the directory map.
├── scanner/ ... Implementation of scanning the target directory and its overview from the file system based on the configuration.
├── scripts/ ... scripts for Dockerfile.
└── version/ ... Versi
```

### `dirmap init`

If you want to change the collection rules, you can create a configuration file ( `.dirmap.yml` ) to change it.

``` console
$ dirmap init
Create .dirmap.yml
$ cat .dirmap.yml
targets:
- file: .dirmap.md
  matcher: markdown
- file: README.md
  matcher: markdown
- file: doc.go
  matcher: godoc
- file: "*.go"
  matcher: godoc
```

## Configuration

### `targets:`

The target files to search for overview document in the directory.

The search for the overview document will be performed on the files specified by `file:` in order.
If the string is matched by `matcher:`, the search in the directory is immediately terminated and the overview document is determined.

``` yaml
targets:
- file: .dirmap.md
  matcher: markdown
- file: README.md
  matcher: markdown
- file: doc.go
  matcher: godoc
- file: '*.go'
  matcher: godoc
```

### `ignores:`

The directories to be excluded.

``` yaml
ignores:
  - dist
  - dist/**/*
```

## Available matcher

### markdown

Get a normal paragraph that is not a heading line without blank lines.

``` yaml
targets:
- file: .dirmap.md
  matcher: markdown
```

### markdownHeading

Get the first heading line.

``` yaml
targets:
- file: .dirmap.md
  matcher: markdownHeading
```

### godoc

Get the text that is retrieved as an overview of the package in the godoc page.

``` yaml
targets:
- file: '*.go'
  matcher: godoc
```

### regexp

If the value of `matcher:` does not match any matcher, it is considered a regular expression.

Then, get the first line matched by the regular expression.

Also, if a capture group is used, only the string that matches the first capture group will be retrieved.

``` yaml
targets:
- file: '*.rb'
  matcher: '^#+\s+(.+)$'
```

## Install

**deb:**

Use [dpkg-i-from-url](https://github.com/k1LoW/dpkg-i-from-url)

``` console
$ export DIRMAP_VERSION=X.X.X
$ curl -L https://git.io/dpkg-i-from-url | bash -s -- https://github.com/k1LoW/dirmap/releases/download/v$DIRMAP_VERSION/dirmap_$DIRMAP_VERSION-1_amd64.deb
```

**RPM:**

``` console
$ export DIRMAP_VERSION=X.X.X
$ yum install https://github.com/k1LoW/dirmap/releases/download/v$DIRMAP_VERSION/dirmap_$DIRMAP_VERSION-1_amd64.rpm
```

**apk:**

Use [apk-add-from-url](https://github.com/k1LoW/apk-add-from-url)

``` console
$ export DIRMAP_VERSION=X.X.X
$ curl -L https://git.io/apk-add-from-url | sh -s -- https://github.com/k1LoW/dirmap/releases/download/v$DIRMAP_VERSION/dirmap_$DIRMAP_VERSION-1_amd64.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/dirmap
```

**go install:**

```console
$ go install github.com/k1LoW/dirmap@vX.X.X
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/dirmap/releases)

**docker:**

```console
$ docker pull ghcr.io/k1low/dirmap:latest
```
