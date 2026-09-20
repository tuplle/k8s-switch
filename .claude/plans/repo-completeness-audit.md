# Audit: what's missing for k8s-switch to be a "complete" public repo

## Context

Audit of what's missing from `k8s-switch` (Go CLI, module
`github.com/tuplle/k8s-switch`) to be a fully complete, professional public
open-source tool repository. Reference/planning document only — findings are
organized by priority so it's easy to decide what to act on later.

Current state at time of writing: 4 commits, no tags, single-file `cmd/` and
`internal/` packages, CI that only builds on push to `main`.

## Tier 1 — Baseline hygiene (what a repo needs before others can rely on it)

- **No automated tests.** Zero `*_test.go` files anywhere; CI doesn't run
  `go test`. `internal/utils.go` (file copy, dir listing, subprocess exec) and
  the flag-filtering logic in `cmd/root.go` are both testable today.
- **CI only triggers on `push: branches: [main]`.** No `pull_request` trigger,
  so external contributions never get automatically linted/built before
  merge. Single OS (`ubuntu-latest`), single Go version — no matrix.
- **No tagged releases.** `git tag --list` is empty. `go install
  github.com/tuplle/k8s-switch@latest` resolves to the latest commit on
  `main`, not a semver release — fine for `go install` but not a real
  "release" a user can pin to, and pkg.go.dev conventions expect tags.
- **No CHANGELOG.md.** No way for a user to see what changed between
  versions.
- **No SECURITY.md.** No documented way to privately report a vulnerability.

## Tier 2 — Release & contribution polish

- **No release automation** (e.g. GoReleaser) — no prebuilt cross-platform
  binaries (macOS/Linux/Windows), no Homebrew tap, no install script. Today
  the only install path is `git clone && make install` (Go toolchain
  required).
- **No README badges** — build status, license, go report card, latest
  release. README currently opens straight into the description.
- **No GitHub issue templates or PR template** (`.github/ISSUE_TEMPLATE/`,
  `PULL_REQUEST_TEMPLATE.md`) despite `CONTRIBUTING.md` inviting PRs/issues.
- **No CODEOWNERS.**
- **No Dependabot/Renovate config** — the two direct deps (`promptui`,
  `cobra`) and four indirect ones won't get automated update PRs;
  `make install-deps` runs `go get -u .` which mutates versions rather than
  pinning them deliberately.

## Tier 3 — Deeper polish (nice-to-have, lower urgency)

- **Linting is `go vet` + `go fmt` only** — no `golangci-lint` config for
  richer static analysis (unused code, shadowing, error-check linting, etc.).
- **No shell completion shipped.** Cobra auto-adds a `completion` subcommand,
  but nothing documents it in the README or ships a generated completion
  script.
- **README gaps**: no demo screenshot/GIF, no "why this exists" / motivation
  section, no comparison to alternatives (e.g. `kubectx`, `kubie`), no
  uninstall instructions.
- **No `.editorconfig` or `.gitattributes`.**
- **No code coverage tooling** (codecov config/badge) — moot until tests
  exist.
- **No package-level Go doc comments** (`doc.go` or `// Package cmd ...`) —
  affects how the module renders on pkg.go.dev. Low impact since the only
  public package (`cmd`) is a CLI wrapper, not a library API most people
  would import.
- **LICENSE copyright line** reads `Copyright 2026 tuplle` — the GitHub
  org/username rather than a legal name; commonly fine for OSS but worth a
  conscious choice rather than a template default.

## Not a gap

`LICENSE.txt` (full genuine Apache 2.0 text), `go.sum` (consistent with
`go.mod`), `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `.gitignore`, and a
working `Makefile` (`build`/`install`/`lint`/`clean`) are all already in
good shape.

## Suggested order if/when the user wants to act

1. Add tests + `go test` in CI + `pull_request` trigger (Tier 1) — cheapest,
   highest confidence payoff.
2. Cut a first semver tag (e.g. `v0.1.0`) once tests exist, add
   CHANGELOG.md and SECURITY.md.
3. Move to Tier 2 (GoReleaser, badges, templates) once the project is ready
   for outside contributors.
4. Tier 3 as ongoing polish.
