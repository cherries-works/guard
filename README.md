# Guard

Dependency security and outdated package analyzer written in Go.

Guard analyzes your project's dependencies across multiple ecosystems, detects known vulnerabilities, and identifies outdated packages.

The simplicity of Guard is intentional. No accounts, or configuration. Point Guard at a project and get a clear overview of its dependencies.

Fast and practical. Keep your dependencies secure and up to date.

## Features

* Dependency analysis
* Vulnerability detection
* Outdated package detection
* Multiple ecosystem support
* Recursive project scanning
* OSV vulnerability database
* Clear CLI output
* Lightweight
* Simple

## Supported Ecosystems

* npm
* Go
* Cargo
* Python
* Maven

## Philosophy

Guard aims to make dependency security simple:

* Fast
* Lightweight
* Self-hosted
* Easy to use
* No configuration
* Clear results
* Developer focused

## Installation

```bash
$ git clone https://github.com/cherries-works/guard.git
$ cd guard
$ go build -o guard ./cmd/guard/
$ ./guard .
```

## License

MIT