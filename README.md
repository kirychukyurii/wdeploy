# wdeploy

<p>
    <a href="https://github.com/kirychukyurii/wdeploy/releases"><img src="https://img.shields.io/github/release/kirychukyurii/wdeploy.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/kirychukyurii/wdeploy?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/kirychukyurii/wdeploy/actions"><img src="https://github.com/kirychukyurii/wdeploy/workflows/golangci-lint/badge.svg" alt="Lint Status"></a>
</p>

Deploy Webitel services on multiple servers

[![asciicast](https://asciinema.org/a/ROgQvKvgF86XRpG7hDeDUcBsg.svg)](https://asciinema.org/a/ROgQvKvgF86XRpG7hDeDUcBsg)

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/kirychukyurii/wdeploy/main/scripts/install.sh | bash
export PATH=$PATH:$HOME/.local/bin
wdeploy --version
```

You can also download a binary from the [releases][releases] page. Packages are
available in Alpine, Debian, and RPM formats. Binaries are available for Linux,
macOS, and Windows.

[releases]: https://github.com/charmbracelet/soft-serve/releases

## Help

```bash
$ wdeploy run --help
Run wdeploy TUI

Usage:
  wdeploy run [flags]

Examples:
wdeploy run --user "testUser" --password "testPassword" --deploy-type custom

Flags:
  -t, --deploy-type string   specify Ansible inventory template type: localhost, custom (default "localhost")
  -h, --help                 help for run
  -i, --inventory string     specify Ansible inventory host path
  -F, --log-format string    log output format: json, console (default "plain")
  -l, --log-level string     log output level: debug, info, warn, error, dpanic, panic, fatal (default "debug")
  -L, --log-path string      log output to this directory (default "./")
  -p, --password string      specify Webitel Repository password
  -u, --user string          specify Webitel Repository user
  -V, --vars string          specify Ansible variables file
```

## Run

```bash
wdeploy run --user "webitel" --password "demo" --deploy-type local --log-level info
```