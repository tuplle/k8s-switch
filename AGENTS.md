# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project overview

`k8s-switch` is a small Go CLI tool that lets a user interactively switch between
multiple kubeconfig files. It lists YAML files from `~/.kube/config.d/`, prompts the
user to pick one via a Bubble Tea / bubbles-list picker, copies the selection over
`~/.kube/config`, and can optionally launch `k9s` or `kubetail` with that config.

Module path: `github.com/tuplle/k8s-switch` (Go 1.25).

## Layout

- `main.go` — entry point, delegates to `cmd.Execute()`.
- `cmd/root.go` — the single Cobra root command: flag definitions, orchestration, and
  thin wrappers (`runK9s`, `runKubetail`) that shell out to external TUIs.
  `cmd/root_test.go` covers it. `cmd/` holds only Cobra command wiring — non-Cobra
  logic (including TUI components) lives in `internal/`, even logic that's only
  used interactively; see `internal/select.go` below.
- `internal/utils.go` — small, dependency-free helpers (`GetFilesFromDir`,
  `FilterByExtension`, `CopyFile`, `RunAnotherTUI`) shared by `cmd`.
  `internal/utils_test.go` covers it.
- `internal/select.go` — the interactive config picker: a Bubble Tea model
  (`charm.land/bubbletea/v2`) wrapping a `charm.land/bubbles/v2/list`. `enter`
  selects; `q`/`esc`/`ctrl+c` cancel — handled explicitly in this model's own
  `Update`, not via the list's built-in quit keybindings (disabled via
  `list.DisableQuitKeybindings()`: bubbles v2.2.1's default `Quit` binding is `v`
  mislabeled "select", which is broken/misleading, so don't re-enable it). Each
  item's description is `"<context> · <server>"`, read from the kubeconfig file
  itself via `describeConfig`/`ReadKubeconfigSummary` (falls back to the file path
  if the file can't be parsed). `internal.SelectConfig` is the entry point
  `cmd/root.go` calls. `internal/select_test.go` covers the model directly (no
  real TTY needed).
- `internal/kubeconfig.go` — `ReadKubeconfigSummary` parses a kubeconfig YAML file
  (via `gopkg.in/yaml.v3`) and returns its context name and cluster server address.
  Assumes exactly one cluster and one context per file (reads only the first entry
  of each; no merging, `current-context` is not consulted) — this assumption is a
  deliberate project convention, not a general kubeconfig-parsing library.
  `internal/kubeconfig_test.go` covers it.
- `bin/` — build output from `make build` (git-ignored, not source).
- `.goreleaser.yaml` — cross-platform release config (see "Releases" below).
- `.github/workflows/build-main.yml` — CI: installs deps, lints, tests, builds on
  push and pull requests targeting `main`.
- `.github/workflows/release.yml` — builds and publishes cross-platform binaries to
  GitHub Releases via GoReleaser when a `v*` tag is pushed.
- `.config/.git-pre-commit` — pre-commit hook source (`go fmt ./...`, `go vet ./...`).

## Build, lint, run

Use the `Makefile` targets rather than raw `go` invocations where one exists:

```bash
make build         # go build -> bin/k8s-switch
make install       # go install github.com/tuplle/k8s-switch
make lint          # go vet ./... && go fmt ./...
make test          # go test ./...
make clean         # remove bin/ and any stray binary
make install-deps  # go mod download (used in CI; does NOT change dependency versions)
make update-deps   # go get -u ./... && go mod tidy (deliberate dependency bump only)
```

`install-deps` must stay non-mutating (`go mod download`, not `go get -u`) — it's part of
the normal build/test/release path and CI, so it must never change `go.mod`/`go.sum`.
Bumping dependency versions is a separate, deliberate action via `make update-deps`,
done in its own commit/PR, not as a side effect of building or releasing.

Run `go test ./...` (or `make test`) and `go fmt ./...` / `go vet ./...` before
committing (mirrors the pre-commit hook and CI's `make lint`/`make test`) — CI will
fail the build otherwise.

## Conventions

- `cmd/` contains only Cobra command definitions and wiring (flags, `Run` orchestration).
  Everything else — helpers, side-effecting logic, and TUI components (prompts,
  Bubble Tea models) — belongs in `internal/`, even when only used interactively from
  one command. `internal/utils.go`/`internal/select.go` functions are small and
  single-purpose — follow that pattern rather than growing one large helper file.
- Exported functions in `internal/` use Go doc comments starting with the function
  name (see `GetFilesFromDir`, `CopyFile`, `RunAnotherTUI`); match that style for new
  exported symbols.
- The root command panics on unexpected errors (e.g. failing to read the config
  directory) rather than wrapping them into a graceful CLI error path — match this
  existing behavior unless asked to change it.
- Flags are defined in `cmd/root.go`'s `init()`; add new flags there and read them
  inside `rootCmd.Run` via `cmd.Flags().Get*`.
- No third-party test or mocking framework is present — stdlib `testing` is the
  expected default if tests are added.

## Releases

`cmd.Version` is `var Version = "dev"` (`cmd/root.go`), not a hardcoded constant —
GoReleaser injects the real version at build time via `-ldflags -X` (see
`.goreleaser.yaml`), so local builds always report `dev`. To cut a release: update
`CHANGELOG.md`, then tag `vX.Y.Z` and push the tag (a git operation left to the human
maintainer, not performed automatically). Pushing the tag triggers
`.github/workflows/release.yml`, which runs GoReleaser to cross-compile for
linux/darwin/windows (amd64+arm64) and publish archives + checksums to GitHub
Releases.

To test the release pipeline locally without tagging or publishing:
`make snapshot` (requires `goreleaser` installed:
`go install github.com/goreleaser/goreleaser/v2@latest`).

## External dependencies at runtime

The tool shells out to external binaries that are not Go dependencies: `k9s` and
`kubetail` must be present on `PATH` for the `--k9s`/`--k9s-only`/`--logs` flags to
work. Don't assume they're installed in a sandboxed/CI environment — code that
exercises `runK9s`/`runKubetail`/`RunAnotherTUI` can't be tested there without them.

## Git / PR notes

- Default branch is `main`; `build-main.yml` runs on push and pull requests targeting
  `main`; `release.yml` runs separately on `v*` tag pushes.
- `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md` document human contributor process —
  worth checking if a change touches project process, licensing, or community docs.
