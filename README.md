<p align="center"><a href="#readme"><img src="https://gh.kaos.st/source-index.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://travis-ci.org/essentialkaos/source-index"><img src="https://travis-ci.org/essentialkaos/source-index.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/source-index"><img src="https://goreportcard.com/badge/github.com/essentialkaos/source-index"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-source-index-master"><img alt="codebeat badge" src="https://codebeat.co/badges/dec317bf-9da2-4d56-ab9b-a31dde545285" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

SourceIndex is a simple utility for generating index page for source code archives. Currently, we use SourceIndex for generating an index for [EK Sources Repository](https://source.kaos.st).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the SourceIndex from scratch, make sure you have a working Go 1.10+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/source-index
```

If you want update SourceIndex to latest stable release, do:

```
go get -u github.com/essentialkaos/source-index
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```
[sudo] source-index --completion=bash 1> /etc/bash_completion.d/source-index
```


ZSH:
```
[sudo] source-index --completion=zsh 1> /usr/share/zsh/site-functions/source-index
```


Fish:
```
[sudo] source-index --completion=fish 1> /usr/share/fish/vendor_completions.d/source-index.fish
```

### Usage

```
Usage: source-index {options} dir

Options

  --output, -o file      Output file (index.html by default)
  --template, -t file    Template (template.tpl by default)
  --no-color, -nc        Disable colors in output
  --help, -h             Show this help message
  --version, -v          Show version

```

### Build Status

| Branch | Status |
|------------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/source-index.svg?branch=master)](https://travis-ci.org/essentialkaos/source-index) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/source-index.svg?branch=develop)](https://travis-ci.org/essentialkaos/source-index) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
