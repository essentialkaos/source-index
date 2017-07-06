# SourceIndex [![Build Status](https://travis-ci.org/essentialkaos/source-index.svg?branch=master)](https://travis-ci.org/essentialkaos/source-index) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/source-index)](https://goreportcard.com/report/github.com/essentialkaos/source-index) [![License](https://gh.kaos.io/ekol.svg)](https://essentialkaos.com/ekol)

SourceIndex is a simple utility for generating index page for source code archives. Currently, we use SourceIndex for generating an index for [EK Sources Repository](https://source.kaos.io).

* [Installation](#installation)
* [Usage](#usage)
* [Build Status](#build-status)
* [Contributing](#contributing)
* [License](#license)

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the SourceIndex from scratch, make sure you have a working Go 1.6+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/source-index
```

If you want update SourceIndex to latest stable release, do:

```
go get -u github.com/essentialkaos/source-index
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

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
