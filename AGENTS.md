# IdeaMapper — Agent Operational Guidelines & Mandatory Protocols

You are working on **IdeaMapper** (`com.antigravity.ideamemoryapp`), located at `/home/aditya/ai_projects/idea_app`.

## 🚨 MANDATORY PROTOCOL: Continuous Living Documentation

Before completing ANY turn or concluding any task in this workspace, you MUST ensure that all living documentation files in `documents/` and `.agent/rules/` are kept strictly synchronized and updated:

### 1. `documents/04_COMMAND_HISTORY.md` (MANDATORY AUTO-UPDATE ON EVERY TURN)
- At the end of every turn where terminal commands, builds, tests, git operations, or script runs occurred:
  - Append a timestamped log entry recording the exact commands executed, arguments/parameters passed, exit code/status, and a concise summary of files modified.
  - Maintain the chronological sequence since project inception. Never leave executed operations unrecorded.

### 2. `documents/03_PROMPTS_AND_DECISIONS.md` (UPDATE ON PROMPTS & DECISIONS)
- Record every significant user request, prompt, or directive.
- Document Architectural Decision Records (ADRs) explaining technical choices, trade-offs, and rationale (e.g., schema refactors, new libraries, UI architecture, anti-greed philosophy).

### 3. `documents/01_APP_DESIGN_AND_ROADMAP.md` (UPDATE ON FEATURE & ROADMAP PROGRESS)
- Whenever a feature is implemented, updated, or scoped, update the V1 MVP Scope or V2/V3 Roadmap sections.
- Mark roadmap checkboxes as completed (`[x]`) as new capabilities are delivered.

### 4. `documents/02_SETUP_AND_INSTALLATION.md` (UPDATE ON DEPENDENCY & SETUP CHANGES)
- If `go.mod` dependencies change, new environment variables are added, or build flags/scripts are modified, update this guide.
- Ensure all prerequisite commands and paths (Go, Fyne, Android SDK API 34, NDK r26d) remain 100% reproducible.

### 5. Core Engineering Rules & Architecture
- Strictly adhere to [`.agent/rules/engineering-standards.md`](./.agent/rules/engineering-standards.md) and [`.agent/rules/app-architecture.md`](./.agent/rules/app-architecture.md):
  - **Database & Concurrency:** Enforce SQLite recursive triggers (`_pragma=recursive_triggers(1)`), defer `rows.Close()` immediately to prevent connection deadlocks, and keep AI/backup I/O off the Fyne UI event thread.
  - **Polymorphic Integrity:** Enforce metadata hierarchy bounds (`Level` 1..6) and cascade deletions.
  - **Testing:** Every change must pass `go test -v ./...`.
  - **Git Hygiene:** Strictly obey GitHub's 100 MB file limit. NEVER commit `*.apk` or `*.aab` files to git. Keep them ignored via `.gitignore`.
