# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project overview

`k8s-switch` is a small Go CLI tool that lets a user interactively switch between
multiple kubeconfig files. It lists YAML files from `~/.kube/config.d/`, prompts the
user to pick one (via `promptui`), copies the selection over `~/.kube/config`, and
can optionally launch `k9s` or `kubetail` with that config.

Module path: `github.com/tuplle/k8s-switch` (Go 1.25).

## Layout

- `main.go` — entry point, delegates to `cmd.Execute()`.
- `cmd/root.go` — the single Cobra root command: flag definitions, the interactive
  selection flow, and thin wrappers (`runK9s`, `runKubetail`) that shell out to
  external TUIs.
- `internal/utils.go` — small, dependency-free helpers (`GetFilesFromDir`,
  `CopyFile`, `RunAnotherTUI`) shared by `cmd`.
- `bin/` — build output from `make build` (git-ignored, not source).
- `.github/workflows/build-main.yml` — CI: installs deps, lints, builds on push to
  `main`.
- `.config/.git-pre-commit` — pre-commit hook source (`go fmt ./...`, `go vet ./...`).

There are no other packages and no test files (`*_test.go`) in the repo currently.

## Build, lint, run

Use the `Makefile` targets rather than raw `go` invocations where one exists:

```bash
make build         # go build -> bin/k8s-switch
make install       # go install github.com/tuplle/k8s-switch
make lint          # go vet ./... && go fmt ./...
make clean         # remove bin/ and any stray binary
make install-deps  # go get -u . (used in CI)
```

There is no test suite yet. If you add one, run it with `go test ./...` and prefer
adding a `test` target to the Makefile alongside it rather than inventing a separate
convention.

Before committing, run `go fmt ./...` and `go vet ./...` (mirrors the pre-commit hook
and CI's `make lint`) — CI will fail the build otherwise.

## Conventions

- Keep `cmd/` focused on CLI wiring (flags, prompts, orchestration) and put reusable,
  side-effecting logic in `internal/`. `internal/utils.go` functions are small and
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

## External dependencies at runtime

The tool shells out to external binaries that are not Go dependencies: `k9s` and
`kubetail` must be present on `PATH` for the `--k9s`/`--k9s-only`/`--logs` flags to
work. Don't assume they're installed in a sandboxed/CI environment — code that
exercises `runK9s`/`runKubetail`/`RunAnotherTUI` can't be tested there without them.

## Git / PR notes

- Default branch is `main`; CI (`build-main.yml`) only triggers on push to `main`,
  there is no PR-triggered workflow currently.
- `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md` document human contributor process —
  worth checking if a change touches project process, licensing, or community docs.
