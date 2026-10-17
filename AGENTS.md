# AGENTS.md — mdbind Repository

Agent guide for AI agents working inside the `mdbind` repository
(module = `github.com/krewire/mdbind`).

`mdbind` is a zero-dependency, standalone Markdown-based book and documentation
site builder written in Go. It operates both as an independent CLI tool (`cmd/mdbind`)
and powers the `book` workload within the Krewire ecosystem (`kiw build`).

## Language

- **Write all stored artifacts in English** — code, docs, commit messages, specs, reports.
- If the human writes in Indonesian or another language, reply in the human's language in chat, but keep all repository files in English.

## Task & Progress Tracking (GitHub Issues SSOT)

To keep progress tracking fully automated and deterministic for both AI agents and human developers:

1. **GitHub Issues as Single Source of Truth:**
   - Active tasks, features, and fixes are tracked as GitHub Issues under `krewire/mdbind`.
   - Every issue body contains a structured checklist with `- [ ]` and `- [x]` checkboxes.

2. **Automated Agent & Developer Progress Workflow:**
   - **Discovering tasks:** Query active issues via CLI/API:
     ```bash
     gh issue list -R krewire/mdbind --state open --json number,title,url
     gh issue view <number> -R krewire/mdbind --json title,body
     ```
   - **Verifying completion:** Run tests:
     ```bash
     go test ./...
     ```
   - **Updating status:** When a task is implemented and verified by tests, check off the checklist item (`- [ ]` → `- [x]`) in the GitHub Issue body. Commits must reference the issue number (`feat(scope): implement feature (ref #<IssueNumber>)`). Close the issue once all tasks are completed.

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
