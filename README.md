<p align="center"><a href="#readme"><img src="https://gh.kaos.st/source-index.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/source-index/ci"><img src="https://kaos.sh/w/source-index/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/source-index"><img src="https://kaos.sh/r/source-index.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/source-index"><img src="https://kaos.sh/b/dec317bf-9da2-4d56-ab9b-a31dde545285.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/source-index/codeql"><img src="https://kaos.sh/w/source-index/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

</br>

SourceIndex is a simple utility for generating index page for source code archives. Currently, we use SourceIndex for generating an index for [EK Sources Repository](https://source.kaos.st).

### Installation

To build the SourceIndex from scratch, make sure you have a working Go 1.17+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/source-index
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo source-index --completion=bash 1> /etc/bash_completion.d/source-index
```


ZSH:
```bash
sudo source-index --completion=zsh 1> /usr/share/zsh/site-functions/source-index
```


Fish:
```bash
sudo source-index --completion=fish 1> /usr/share/fish/vendor_completions.d/source-index.fish
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
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/source-index/ci.svg?branch=master)](https://kaos.sh/w/source-index/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/source-index/ci.svg?branch=master)](https://kaos.sh/w/source-index/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
