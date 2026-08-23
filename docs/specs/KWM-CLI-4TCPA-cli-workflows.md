# Specification — mdbind CLI & Workflows (Superseded)

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | KWM-4TCPA                                   |
| Title       | mdbind CLI & Workflows                     |
| Status      | Superseded                                  |
| Date        | 2026-08-18                                  |
| Author      | Krewire Contributors                         |
| Domain      | Site Builders — CLI                        |

## 1. Context

The mdbind CLI (now superseded) dogfoods the Krewire Framework's `cli` package and
configuration conventions (`MDBIND_*` environment variables, flag > env > default
precedence). It exposed the builder as commands: `build`, `init`, and `serve`.

**Note:** The CLI commands have been moved to the `krewire` devtool. Use
`krewire build`, `krewire serve`, and `krewire init` instead. The `book` library
remains the public API for building and serving sites programmatically.

## 2. Problem Statement

Without a uniform CLI, every team builds and previews its website differently:
bespoke build scripts, ad-hoc preview servers, and manual scaffolding. The
mdbind CLI normalized that experience — one command to build, one to scaffold,
one to preview — with canonical exit codes and structured diagnostics.

## 3. Goals

- G1 — Provide `build`, `init`, and `serve` commands on the framework `cli` (now in `krewire`).
- G2 — Follow the ecosystem configuration convention (`MDBIND_*` env, flag > env > default).
- G3 — Resolve every outcome to a canonical exit code (0/1/2).
- G4 — Keep data on stdout and diagnostics on stderr.

## 4. Non-Goals

- NG1 — Watch mode, hot reload, or incremental builds in this phase.
- NG2 — Deployment or hosting from the CLI.
- NG3 — Plugin/extension loading.

## 5. Requirements

### 5.1 Build

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-CL-001 | `build` renders the manuscript into a static website.             | Must     |
| MDB-CL-002 | Support `--input`, `--output`, `--title`, `--author`, `--base`, `--nav`, `--footer`, and `--theme` flags. | Must |
| MDB-CL-003 | Read the same settings from `MDBIND_*` environment variables.     | Should   |
| MDB-CL-004 | Default input to `.`, output to `site`.                           | Must     |
| MDB-CL-005 | List created files deterministically on stdout.                   | Must     |

### 5.2 Init

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-CL-006 | `init` scaffolds a sample manuscript in a target directory.       | Must     |
| MDB-CL-007 | Refuse to scaffold into a non-empty directory.                    | Must     |

### 5.3 Serve

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-CL-008 | `serve` previews the manuscript over HTTP through the web router. | Must     |
| MDB-CL-009 | Support `--addr` to choose the listen address.                    | Must     |

### 5.4 Behavior

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-CL-010 | Invalid usage resolves to the usage exit code (2).                | Must     |
| MDB-CL-011 | Build failures resolve to the failure exit code (1).              | Must     |
| MDB-CL-012 | Success resolves to the success exit code (0).                    | Must     |
| MDB-CL-013 | All commands are registered through the framework `cli` package.  | Must     |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Determinism.** Identical invocations produce identical output.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — `mdbind init` then `mdbind build` produces a complete site with exit 0.
- S2 — `mdbind serve` answers chapter requests with 200 and unknown paths with 404.
- S3 — Bad flags and missing input directories exit 2 and 1 respectively.
- S4 — `build` output matches the `book.Build` library output byte-for-byte.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [KWM-FX9H2](./KWM-BUILDER-FX9H2-mdbind-site-builder.md)  | mdbind Site Builder — Initial Specification |
| [KWN-Z0VFC](https://github.com/krewire/krewire/blob/main/docs/specs/KWN-DEVTOOL-Z0VFC-krewire-devtool.md) | Krewire Devtool — Initial Specification |
| [KWF-M07QS](https://github.com/krewire/framework/blob/main/docs/specs/KWF-WEB-M07QS-krewire-web-framework.md) | Krewire Web Framework |
| [KWF-FGNZ9](https://github.com/krewire/framework/blob/main/docs/specs/KWF-CLI-FGNZ9-cli-configuration.md) | CLI Configuration |
| [KWF-KAKQL](https://github.com/krewire/framework/blob/main/docs/specs/KWF-CLI-KAKQL-cli-errors-diagnostics.md) | CLI Errors & Diagnostics |

## 9. References

- [KWL-W0J2X](https://github.com/krewire/libs/blob/main/docs/specs/KWL-CORE-W0J2X-errors-exit-codes.md) — Core Errors & Exit Codes.
- [KWM-FX9H2](./KWM-BUILDER-FX9H2-mdbind-site-builder.md) — the builder backing every command.