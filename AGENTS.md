# AGENTS.md — mdbind Repository

Agent guide for AI agents working inside the `mdbind` repository
(module = `github.com/krewire/mdbind`).

`mdbind` is a zero-dependency, standalone Markdown-based book and documentation
site builder written in Go. It operates both as an independent CLI tool (`cmd/mdbind`)
and powers the `book` workload within the Krewire ecosystem (`kiw build`).

## Language

- **Write all stored artifacts in English** — code, docs, commit messages, specs, reports.
- If the human writes in Indonesian or another language, reply in the human's language in chat, but keep all repository files in English.

## Specs & Progress Tracking (GitHub Issues SSOT)

To keep progress tracking fully automated, resilient against markdown parsing drift, and deterministic for both AI agents and human developers:

1. **GitHub Issues as Single Source of Truth:**
   - Active specifications and their implementation progress are tracked as GitHub Issues under `krewire/mdbind` with label `spec`.
   - Title format: `[Spec: KWM-<Code>] <Title>` (e.g. `[Spec: KWM-FX9H2] mdbind Site Builder`, `[Spec: KWM-4TCPA] mdbind Standalone CLI & Workflows`).
   - Every issue body contains the structured requirement checklist with `- [ ]` and `- [x]` checkboxes citing `ReqID`, priority, and scope (e.g. `- [ ] MDB-CL-001 [Must] ...`).
   - `docs/specs/index.md` links directly to these GitHub Issues.

2. **Spec-Driven Traceability in Code:**
   - ID prefix: `KWM-*`.
   - Traceability: test files must cite `// Tests for KWM-<Code>` and individual tests cite requirement IDs (e.g. `MDB-CL-001`, `MDB-BD-001`).

3. **Automated Agent & Developer Progress Workflow:**
   - **Discovering tasks:** Query active specs via CLI/API:
     ```bash
     gh issue list -R krewire/mdbind --label spec --json number,title,url
     gh issue view <number> -R krewire/mdbind --json title,body
     ```
   - **Verifying completion:** Run tests:
     ```bash
     go test ./...
     ```
   - **Updating status:** When a requirement is implemented and verified by tests, check off the requirement box (`- [ ]` → `- [x]`) in the GitHub Issue body. Commits must reference the issue number and requirement ID (`feat(scope): implement <ReqID> (ref #<IssueNumber>, req <ReqID>)`). Close the issue once all requirements are completed.

## Quality Gates

Before finishing any task, run:
```bash
gofmt -l .        # must print nothing
go vet ./...
go test ./...
go build ./cmd/mdbind
```

## Commit Discipline

- **Conventional Commits**: `feat|fix|docs|test|refactor|perf|security|chore(scope): desc`.
- Atomic commits (≤~300 LOC).
- Run `git` inside `mdbind/`.
