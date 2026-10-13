# Specifications Index — mdbind

This directory holds the formal specifications for mdbind.

Ordered by **impact-to-effort** (high impact, low effort first) and **dependency chain** (foundations first).

| SpecID    | Title                                      | Status | Impl Status | Depends On |
| --------- | ------------------------------------------ | ------ | ----------- | ---------- |
| [KWM-FX9H2](https://github.com/krewire/mdbind/issues/1) | mdbind Site Builder — Initial Specification | Draft | Shipped | — |
| [KWM-4TCPA](https://github.com/krewire/mdbind/issues/2)      | mdbind CLI & Workflows                    | Draft | Superseded (standalone CLI → `kiw build` book mode) | KWM-FX9H2 |

## Conventions

- Each specification is stored as a single Markdown file named `{ProjectId}-{Scope}-{SpecID}-{slug}.md`.
- SpecIDs are unique, random 5-character alphanumeric codes (e.g., `KWM-FX9H2`).
- New specifications must be added to this index when created.
- Ordering: impact-to-effort (high impact, low effort first), then dependency chain (foundations first).