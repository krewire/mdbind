# Specifications Index — mdbind

This directory holds the formal specifications for mdbind.

Ordered by **impact-to-effort** (high impact, low effort first) and **dependency chain** (foundations first).

| SpecID    | Title                                      | Status | Depends On |
| --------- | ------------------------------------------ | ------ | ---------- |
| [KWM-FX9H2](./KWM-BUILDER-FX9H2-mdbind-site-builder.md) | mdbind Site Builder — Initial Specification | Draft | — |
| [KWM-4TCPA](./KWM-CLI-4TCPA-cli-workflows.md)      | mdbind CLI & Workflows                    | Draft | KWM-FX9H2 |

## Conventions

- Each specification is stored as a single Markdown file named `{ProjectId}-{Scope}-{SpecID}-{slug}.md`.
- SpecIDs are unique, random 5-character alphanumeric codes (e.g., `KWM-FX9H2`).
- New specifications must be added to this index when created.
- Ordering: impact-to-effort (high impact, low effort first), then dependency chain (foundations first).