# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.0.0] - 2026-09-20

### Added

- Interactive selection of kubeconfig files from `~/.kube/config.d/`, with the
  chosen file copied over `~/.kube/config`.
- `--k9s` / `-9` flag to launch [k9s](https://k9scli.io/) with the selected
  config after switching.
- `--k9s-only` flag to launch k9s with the selected config without touching
  `~/.kube/config`.
- `--logs` flag to launch `kubetail serve` with the selected config.
- `--dir` / `-d` flag to override the kubeconfig source directory.
- `--version` flag.
- Automated test suite (`go test ./...`) covering `internal` and `cmd`.
- CI now runs on pull requests in addition to pushes to `main`, and runs the
  test suite before building.

[1.0.0]: https://github.com/tuplle/k8s-switch/releases/tag/v1.0.0
