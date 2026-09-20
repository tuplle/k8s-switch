# Audit: what's missing for k8s-switch to be a "complete" public repo

## Context

Audit of what's missing from `k8s-switch` (Go CLI, module
`github.com/tuplle/k8s-switch`) to be a fully complete, professional public
open-source tool repository. Findings are organized by priority. Updated to
track what's actually been implemented since the original audit.

Original state audited: 4 commits, no tags, single-file `cmd/` and
`internal/` packages, CI that only builds on push to `main`.

Current state: Tier 1 is complete and released as `v1.0.0` (tag pushed,
matches `origin/main`). Part of Tier 2 is done. A follow-up security audit
also found and fixed several issues not in the original list (see bottom).

## Tier 1 — Baseline hygiene — ✅ DONE (released as v1.0.0)

- [x] **Automated tests** — `internal/utils_test.go` and `cmd/root_test.go`
  added (11 tests), covering file listing, extension filtering, config
  copying, and subprocess execution. The inline YAML-filter logic in
  `cmd/root.go` was extracted into `internal.FilterByExtension` to make it
  testable.
- [x] **CI trigger** — `build-main.yml` now also runs on `pull_request`
  against `main` (previously push-to-main only), and runs `make test`
  before `make build`. (Still single OS/Go-version, no matrix — left as-is,
  not required for baseline.)
- [x] **Tagged release** — `v1.0.0` tag created and pushed; `rootCmd.Version`
  set to `"1.0.0"`, enabling `k8s-switch --version`.
- [x] **CHANGELOG.md** — added, Keep a Changelog format.
- [x] **SECURITY.md** — added, points to GitHub private vulnerability
  reporting.

## Tier 2 — Release & contribution polish

- [ ] **No release automation** (e.g. GoReleaser) — still no prebuilt
  cross-platform binaries, Homebrew tap, or install script. Only install
  path today is `git clone && make install`.
- [x] **README badges** — build status, pkg.go.dev, latest release, and
  license badges added. (Go Report Card badge was proposed but removed by
  the user.)
- [x] **GitHub issue templates and PR template** — added
  `.github/ISSUE_TEMPLATE/bug_report.md`, `feature_request.md`,
  `config.yml` (links to CONTRIBUTING.md), and
  `.github/PULL_REQUEST_TEMPLATE.md`.
- [ ] **No CODEOWNERS.**
- [~] **Dependabot/Renovate config** — still not added, so dependency
  updates aren't automated via PRs. Partially addressed differently:
  `make install-deps` no longer mutates versions (now `go mod download`);
  a separate, deliberate `make update-deps` target
  (`go get -u ./... && go mod tidy`) was added for intentional bumps.

## Tier 3 — Deeper polish (nice-to-have, lower urgency) — not started

- [ ] Linting is still `go vet` + `go fmt` only — no `golangci-lint` config.
- [ ] No shell completion shipped/documented (Cobra's built-in `completion`
  subcommand).
- [ ] README still has no demo screenshot/GIF, "why this exists" section,
  comparison to alternatives (`kubectx`, `kubie`), or uninstall
  instructions.
- [ ] No `.editorconfig` or `.gitattributes`.
- [ ] No code coverage tooling (codecov config/badge) — tests now exist, so
  this is unblocked if wanted.
- [ ] No package-level Go doc comments (`doc.go`).
- [ ] LICENSE copyright line still reads `Copyright 2026 tuplle`.

## Not a gap

`LICENSE.txt`, `go.sum`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`,
`.gitignore`, and a working `Makefile` were already in good shape.

## Security audit — ✅ DONE (found separately, not in original tiers)

A follow-up security review (govulncheck + manual review) found and fixed:

- [x] **`~/.kube/config` written with world-readable permissions** —
  `CopyFile` now writes via `os.CreateTemp` (mode `0600`) instead of
  `os.Create` (was `0644`).
- [x] **Non-atomic config write** — `CopyFile` now writes to a temp file in
  the destination directory, `Sync`s, then `os.Rename`s into place, so an
  interruption never leaves a partially written `~/.kube/config`.
- [x] **Symlinks followed without validation** — `GetFilesFromDir` now
  filters on `entry.Type().IsRegular()` instead of `!entry.IsDir()`, so
  symlinks in the config directory are excluded rather than followed.
- [x] **CI actions pinned by mutable tag** — `actions/checkout@v6` and
  `actions/setup-go@v6` pinned to their resolved commit SHAs.
- [x] **No explicit CI token permissions** — added
  `permissions: contents: read` to `build-main.yml`.
- [x] Dependency vulnerability scan (`govulncheck`) — clean, no known CVEs.

## Suggested order for what's left

1. Tier 2 remainder: CODEOWNERS, Dependabot config, release automation
   (GoReleaser) once ready for outside contributors/prebuilt binaries.
2. Tier 3 as ongoing polish.
