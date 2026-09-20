# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Cross-platform release pipeline: pushing a `v*` tag now builds
  linux/darwin/windows (amd64+arm64) binaries via GoReleaser and publishes
  them as GitHub Release archives with checksums, so users can install a
  prebuilt binary without a Go toolchain.
- `--version` now reports the exact released version on GoReleaser-built
  binaries (injected at build time); local `make build` reports `dev`.

### Changed

- The interactive config picker now uses [Bubble Tea](https://github.com/charmbracelet/bubbletea)
  v2 and `bubbles/list` v2 instead of `promptui`. Keybindings change
  slightly: press `/` to fuzzy-filter by name, and `q`/`esc`/`ctrl+c` (in
  addition to `ctrl+c` alone before) all cancel the selection.

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
