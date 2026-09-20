# Audit: what's missing for k8s-switch to be a "complete" public repo

## Context

Audit of what's missing from `k8s-switch` (Go CLI, module
`github.com/tuplle/k8s-switch`) to be a fully complete, professional public
tool repository. Findings are organized by priority. Updated to track what's
actually been implemented since the original audit.

**Maintainer stance (stated 2026-09-21): not accepting public contributions.**
The repo stays public/usable by anyone, but this is a solo-maintained
project — no external PRs are being solicited. This changes the priority of
a few Tier 2 items below (see "Contribution-related items" note).

Original state audited: 4 commits, no tags, single-file `cmd/` and
`internal/` packages, CI that only builds on push to `main`.

Current state: Tier 1 is complete and released as `v1.0.0`. Tier 2's release
automation (GoReleaser, cross-platform binaries, rpm/deb packages) is now
done. A security audit found and fixed several issues not in the original
list (see below). Since v1.0.0, the interactive picker was migrated from
`promptui` to Bubble Tea v2 + `bubbles/list` v2, and now shows each
kubeconfig's context name + server address instead of its file path. A
`v1.1.0` release is staged in `CHANGELOG.md` but not yet tagged (tagging is
the user's own gitops step, not done by the agent).

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
- [x] **Tagged release** — `v1.0.0` tag created and pushed; `cmd.Version`
  enables `k8s-switch --version` (now build-time injected via GoReleaser
  `-ldflags`, see Tier 2 release automation).
- [x] **CHANGELOG.md** — added, Keep a Changelog format.
- [x] **SECURITY.md** — added, points to GitHub private vulnerability
  reporting.

## Tier 2 — Release & contribution polish

- [x] **Release automation** — `.goreleaser.yaml` + `.github/workflows/release.yml`
  (triggers on `v*` tag push): cross-compiles linux/darwin/windows
  (amd64+arm64), publishes `.tar.gz`/`.zip` archives + `checksums.txt` to
  GitHub Releases, and now also builds `.rpm`/`.deb` packages (via the
  bundled `nfpm`) for the linux builds. `cmd.Version` changed from a
  hardcoded const to a build-time-injected var to support this.
  **Not done, by explicit choice when scoped:** Homebrew tap, install
  script, code signing.
  **New, discussed but not started:** a hosted dnf/apt repo (e.g. Fedora
  Copr) so users get `dnf install k8s-switch` by name with update support,
  rather than downloading and locally installing the `.rpm`/`.deb` from
  each release. Needs an external Copr account + a `.spec` file — action
  requires the user's decision/account, not something the agent can set up
  unilaterally.
- [x] **README badges** — build status, pkg.go.dev, latest release, and
  license badges added. (Go Report Card badge was proposed but removed by
  the user.)
- [x] **GitHub issue templates and PR template** — added
  `.github/ISSUE_TEMPLATE/bug_report.md`, `feature_request.md`,
  `config.yml` (links to CONTRIBUTING.md), and
  `.github/PULL_REQUEST_TEMPLATE.md`.
- [~] **Dependabot/Renovate config** — still not added, so dependency
  updates aren't automated via PRs. Partially addressed differently:
  `make install-deps` no longer mutates versions (now `go mod download`);
  a separate, deliberate `make update-deps` target
  (`go get -u ./... && go mod tidy`) was added for intentional bumps. Not
  affected by the no-public-contributions stance — this is about the
  maintainer's own dependency hygiene, not external PRs.

### Contribution-related items (re-evaluated given "not accepting public contributions")

- **CODEOWNERS — not needed.** This file exists to require specific
  reviewers on PRs; with a solo maintainer not soliciting external
  contributions, there's no one else to route reviews to. Removed from the
  open punch list.
- [x] **`CONTRIBUTING.md` wording softened** — now invites bug reports and
  feature requests, and its "Pull Requests" section explicitly says the
  project is solo-maintained and not accepting PRs (open an issue instead).
  The "Development Setup"/"Project Structure"/"Style Guide" sections were
  removed since they only existed to help prepare a PR.
  `.github/ISSUE_TEMPLATE/config.yml`'s contact-link description was also
  updated (dropped its "or PR" mention).
- **Issue templates/PR template — kept as-is, still useful.** Bug reports
  from users are independent of whether code contributions are accepted;
  the PR template costs nothing to leave in place even if unused.

## Tier 3 — Deeper polish (nice-to-have, lower urgency) — not started

- [ ] Linting is still `go vet` + `go fmt` only — no `golangci-lint` config.
- [x] Shell completion — documented in `README.md` (`## Shell Completion`),
  using Cobra's built-in `completion` subcommand; no static scripts shipped
  in packages/archives.
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
`.gitignore`, and a working `Makefile` were already in
good shape structurally.

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
  `actions/setup-go@v6` (and later `goreleaser-action`) pinned to their
  resolved commit SHAs.
- [x] **No explicit CI token permissions** — `contents: read` on
  `build-main.yml`, `contents: write` (needed to publish releases) on
  `release.yml`.
- [x] Dependency vulnerability scan (`govulncheck`) — clean, no known CVEs
  (at time of scan).

## Post-v1.0.0 feature work (not from the original audit)

- Migrated the interactive picker from `promptui` to Bubble Tea v2 +
  `bubbles/list` v2. Caught and worked around a real upstream bug in
  `bubbles@v2.2.1`'s default keymap (`Quit` bound to `v`, mislabeled
  "select") via explicit key handling instead of relying on the list's
  built-in quit bindings.
- Moved the picker (`select.go`) from `cmd/` to `internal/` — `cmd/` is
  Cobra-wiring-only by project convention now (documented in `AGENTS.md`).
- Picker items now show `"<context> · <server>"` (parsed from the
  kubeconfig YAML via a new `internal/kubeconfig.go`, `gopkg.in/yaml.v3`)
  instead of the file path, assuming one cluster/context per file; falls
  back to the path if a file can't be parsed.
- README's flags table was corrected (`--logs`, `--dir`/`-d`, `--version`
  were missing).

## Suggested order for what's left

1. Tier 2 remainder: Dependabot config; decide on Fedora Copr (or another
   hosted repo) if true `dnf install k8s-switch`-by-name matters enough to
   justify the setup.
2. Tier 3 as ongoing polish.
