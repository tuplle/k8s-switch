# k8s-switch

A lightweight CLI tool to quickly switch between multiple Kubernetes configurations and optionally
launch [k9s](https://k9scli.io/).

## Features

- 📂 **Interactive Selection**: List and select kubeconfigs from your `~/.kube/config.d/` directory.
- 🔄 **Auto-Update**: Automatically overwrites your active `~/.kube/config` with the selected file.
- 🐕 **k9s Integration**: Launch k9s immediately using your newly selected configuration.
- ⚡ **Fast & Simple**: Built with Go, Cobra, and Promptui for a smooth terminal experience.

## Installation

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
| `--verbose`  | `-v`      | Enable verbose output                                                                        |
| `--help`     | `-h`      | Show help message                                                                            |

### Examples

**Switch config and open k9s:**

```shell script
k8s-switch -9
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.txt](LICENSE.txt) file for details.
