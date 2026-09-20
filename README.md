# k8s-switch

[![Build Main](https://github.com/tuplle/k8s-switch/actions/workflows/build-main.yml/badge.svg)](https://github.com/tuplle/k8s-switch/actions/workflows/build-main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tuplle/k8s-switch.svg)](https://pkg.go.dev/github.com/tuplle/k8s-switch)
[![Latest Release](https://img.shields.io/github/v/release/tuplle/k8s-switch)](https://github.com/tuplle/k8s-switch/releases)
[![License](https://img.shields.io/github/license/tuplle/k8s-switch)](LICENSE.txt)

A lightweight CLI tool to quickly switch between multiple Kubernetes configurations and optionally
launch [k9s](https://k9scli.io/).

## Features

- 📂 **Interactive Selection**: List and select kubeconfigs from your `~/.kube/config.d/` directory.
- 🔄 **Auto-Update**: Automatically overwrites your active `~/.kube/config` with the selected file.
- 🐕 **k9s Integration**: Launch k9s immediately using your newly selected configuration.
- ⚡ **Fast & Simple**: Built with Go, Cobra, and Bubble Tea for a smooth terminal experience.

## Installation

### Prebuilt Binaries

Download the archive for your OS/architecture from the
[Releases page](https://github.com/tuplle/k8s-switch/releases), extract it,
and put the `k8s-switch` binary on your `PATH`. No Go toolchain required.

### From Source

Ensure you have Go installed (version 1.25 or later).

```bash
git clone https://github.com/tuplle/k8s-switch.git
cd k8s-switch
make install
```

This will install the `k8s-switch` binary into your `$GOPATH/bin`.

## Setup

Place your various Kubernetes configuration files into the following directory:
`~/.kube/config.d/`

Example:

```plain text
~/.kube/conf.d/
├── prod-cluster.yaml
├── staging-cluster.yaml
└── dev-local.conf
```

## Usage

Simply run the command to start the interactive prompt:

```shell script
k8s-switch
```

### Options

| Flag         | Shorthand | Description                                                                                  |
|:-------------|:----------|:---------------------------------------------------------------------------------------------|
| `--k9s`      | `-9`      | Launch k9s with the selected configuration                                                   |
| `--k9s-only` |           | Only launch k9s with the selected configuration. It does not copy the config to .kube/config |
| `--logs`     |           | Open logs in the default browser with kubetail, using the selected configuration             |
| `--dir`      | `-d`      | Path to the directory containing kubeconfig files (default `~/.kube/config.d/`)              |
| `--verbose`  | `-v`      | Enable verbose output                                                                        |
| `--version`  |           | Show the k8s-switch version                                                                  |
| `--help`     | `-h`      | Show help message                                                                            |

### Examples

**Switch config and open k9s:**

```shell script
k8s-switch -9
```

## Shell Completion

k8s-switch supports shell completion for bash, zsh, fish, and PowerShell (via
[Cobra](https://github.com/spf13/cobra)).

**Bash:**

```bash
source <(k8s-switch completion bash)
# or, to persist across sessions:
k8s-switch completion bash > /etc/bash_completion.d/k8s-switch
```

**Zsh:**

```bash
source <(k8s-switch completion zsh)
# or, to persist across sessions:
k8s-switch completion zsh > "${fpath[1]}/_k8s-switch"
```

**Fish:**

```bash
k8s-switch completion fish | source
# or, to persist across sessions:
k8s-switch completion fish > ~/.config/fish/completions/k8s-switch.fish
```

**PowerShell:**

```powershell
k8s-switch completion powershell | Out-String | Invoke-Expression
```

Run `k8s-switch completion --help` for more details.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.txt](LICENSE.txt) file for details.
