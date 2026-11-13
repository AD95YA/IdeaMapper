# IdeaMapper — Engineering & Agent Rules

## 1. Concurrency, Database Discipline & UI Responsiveness
- **SQLite Connection & Cursor Discipline:** Always ensure `rows.Close()` is deferred or executed immediately after database queries. In SQLite configurations with single-writer connections (`MaxOpenConns=1`), never execute nested queries against the same database handle while an active row cursor is open; read all rows into memory slices first to prevent connection deadlocks.
- **Recursive Triggers & Foreign Keys:** Every SQLite connection string must enforce `_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)`. Never disable or bypass cascading triggers.
- **Non-Blocking UI Event Loop:** Never perform heavy I/O, network requests, database exports, or Gemini AI API calls directly on Fyne's main event thread. Execute long-running tasks inside background goroutines (`go func()`) and marshal state back to the UI cleanly.
- **Pure Go & CGo Independence:** Maintain pure Go compilation for desktop binaries using `modernc.org/sqlite` to ensure fast builds and prevent host C library version mismatches. Reserve CGo strictly for mobile cross-compilation via `fyne package -os android`.

## 2. Domain Integrity, Validation & Anti-Exploit Principles
- **Descriptive Naming Standards:** Use clear, unambiguous, domain-driven names for structs, repositories, and UI constructors (e.g., `GetMetaDataForParent`, `BuildInterLinksSection`, `ShowCreateIdeaModal`). Avoid cryptic abbreviations outside standard iteration counters (`i`, `idx`).
- **Polymorphic Invariant Enforcement:** Strictly validate polymorphic relationships prior to persistence:
  - Enforce `ParentType` to be one of `"Idea"`, `"Memory"`, `"InterLink"`, or `"MetaData"`.
  - Enforce metadata hierarchy bounds (`1 <= Level <= 6`). Reject invalid nesting depths.
  - Require valid UUID identifiers on all created entities using the crypto-random UUID generator.
- **Cascade Cleanliness:** Deleting an Idea, Memory, or InterLink must recursively delete all associated sub-thoughts and file attachments. Never leave orphaned records in `metaData` or `attachments`.
- **No Shortcuts or Fake Successes:** Never bypass validation checks, hardcode fake success returns, or disable tests to force a build to pass.

## 3. Robust Error Handling & Defensive Logging
- **Contextual Error Wrapping:** Wrap all database and external API failures with descriptive context using `fmt.Errorf("action context: %w", err)` to maintain traceability.
- **User-Facing vs. Developer Diagnostics:** Surface user-friendly error dialogs in the UI using `dialog.ShowError(err, w)` while logging detailed diagnostic information to standard logging channels (`log.Printf`).
- **AI & Cloud Resilience:** Gracefully handle Gemini API network failures, invalid or missing API keys, rate limits, and context timeouts without crashing the application. Provide clear visual cues in the Settings and generation dialogs.

## 4. Comprehensive Edge-Case Unit Testing
- Every core package (`db`, `api`, `ai`, `models`) must maintain passing unit tests covering both standard flows and extreme edge cases:
  - **Database & Repositories (`db`):** In-memory SQLite (`:memory:`), multi-level metadata recursion (Levels 1 through 6), cascade deletion triggers on parent removal, circular InterLink creation, and concurrent read access.
  - **Backup & Serialization (`api`):** Full round-trip JSON backup generation and restoration, handling empty databases, restoring nested metadata trees, handling malformed/corrupted JSON payloads, and preventing ID duplication.
  - **AI Integration (`ai`):** Missing or empty API key detection, environment fallback to `GEMINI_API_KEY`, context cancellation, and empty prompt handling.
  - **Domain Models (`models`):** RFC4122 v4 UUID format validation, timestamp preservation, and struct pointer safety.

## 5. Living Documentation & Auto-Updated Command History (`documents/`)
- Maintain a `documents/` folder in the workspace root containing:
  - `documents/01_APP_DESIGN_AND_ROADMAP.md` (or `documents/01_GAME_DESIGN_AND_ROADMAP.md`) — Full app concept, V1 MVP scope (Ideas, Memories, recursive MetaData, InterLinks, Gemini AI, JSON backup), and detailed V2/V3 backlog.
  - `documents/02_SETUP_AND_INSTALLATION.md` — WSL2, Go, Fyne graphics dependencies, Android SDK (API 34) & NDK r26d setup steps, installed tool versions, and build commands (`go run`, `go build`, `./build_apk.sh`).
  - `documents/03_PROMPTS_AND_DECISIONS.md` — Archive of design prompts, title/package decisions, architectural refactors (polymorphic metadata, recursive SQLite cascade triggers, ModernC pure Go driver), and the Anti-Greed Mindly alternative manifesto.
  - `documents/04_COMMAND_HISTORY.md` — **MANDATORY AUTO-UPDATE:** After every agent turn or terminal execution, append a timestamped entry logging the exact shell commands executed, flags/parameters passed, exit status, and brief summary of files modified since the inception of the project.
- Also maintain [`docs/COMMANDS_AND_CHANGES.md`](file:///home/aditya/ai_projects/idea_app/docs/COMMANDS_AND_CHANGES.md) and [`README.md`](file:///home/aditya/ai_projects/idea_app/README.md) with synchronized historical and user-facing records.

## 6. GitHub Repository & Android Packaging Workflow
- **Remote Repository:** `https://github.com/abhokarikar/IdeaMapper.git` (active branch: `main`).
- **Git Large File Size Limit (100 MB):** GitHub strictly enforces a 100 MB per-file limit. Never commit or push pre-built APK files (`IdeaMapper.apk`, `*.apk`, `*.aab`) or compiled binary executables (`idea_app`) to git. Keep all binary artifacts excluded via `.gitignore`.
- **Android APK Build Pipeline:** Use the included `./build_apk.sh` script with Android SDK (API 34) and NDK r26d (`26.3.11579264`). Ensure the Fyne CLI produces valid signed release/debug packages without manual source patches.
- **Commit Cadence:** Stage and commit logical increments using conventional commit prefixes (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `build:`), keeping commit messages descriptive and concise.

