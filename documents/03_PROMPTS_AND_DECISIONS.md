# IdeaMapper — Design Prompts, Decisions & Architecture Log

## 1. Project Genesis & Prompt Archive

### Prompt 1: Initial Inception (Phase 0)
> *"App to help you map all the thoughts and ideas in your brain."*
- **Outcome:** Initial repository layout established on Mar 28, 2026. SQLite database schema laid out for `ideas`, `memories`, `interlinks`, and initial notes stub.

### Prompt 2: Metadata Hierarchy Refactor (Phase 0 / Commit `72ae945`)
> *"renamed notes to metadata"*
- **Outcome:** Initiated the conceptual shift from flat one-dimensional "Notes" to a recursive, multi-level polymorphic "MetaData" tree supporting parent-child sub-thoughts up to 6 levels deep. Left database repositories and tests temporarily unaligned.

### Prompt 3: Compilation Fixes & Core Completion (Phase 1–3)
> *"Fix compilation errors and update unit tests for metadata, links, and triggers. Implement full detail views, hierarchical metadata, interlinks mind map, Gemini integration, and backup UI."*
- **Outcome:** Resolved connection pool deadlocks, enabled SQLite recursive triggers, updated GenAI client signatures for Gemini 2.5 Flash, implemented 4-tab Fyne bottom navigation, and built complete test suites.

### Prompt 4: Android Cross-Compilation & Packaging (Phase 5–7)
> *"Generate Android APK for physical device usage."*
- **Outcome:** Probed system toolchains, installed Android SDK command-line tools, platform `android-34`, build-tools `34.0.0`, and NDK `26.3.11579264` (r26d LTS). Created automated packaging script `build_apk.sh` and generated 174 MB universal multi-ABI `IdeaMapper.apk`.

### Prompt 5: Workspace Relocation (Phase 8)
> *"Migrate project from scratch directory to permanent development workspace."*
- **Outcome:** Transferred full project to `/home/aditya/ai_projects/idea_app`, preserving 100% of git history, build artifacts, test suites, and documentation.

### Prompt 6: README Documentation & Origin Story (Phase 9)
> *"create README.md file and give all details about the app and how to best use it... mention that this was specifically created as the Mindly app on android was unceremoniously taken down by its creators... Write this in a more sarcastic and creative manner."*
- **Outcome:** Authored comprehensive README with 5-step knowledge workflow guide, Antigravity IDE attribution, and the Anti-Greed Manifesto critiquing Mindly's corporate rug-pull. Resolved GitHub 100 MB file limit via `.gitignore` and pushed to GitHub.

### Prompt 7: Living Documentation & Auto-Update Enforcement (Phase 10)
> *"ensure all these files keep getting auto updated with relevant info for any future work"*
- **Outcome:** Built a multi-layered automated enforcement architecture: authored root `AGENTS.md` and `GEMINI.md` workspace rules, configured Antigravity lifecycle hooks in `.agent/hooks.json` via `scripts/hooks/pre_invocation.sh`, and provided `scripts/log_command.sh` to ensure all future agent turns and terminal operations continuously update `documents/`.

---

## 2. Key Architectural Decisions & Trade-Offs

### ADR-01: Pure Go SQLite (`modernc.org/sqlite`) vs. CGo (`mattn/go-sqlite3`)
- **Context:** SQLite was needed for local embedded storage. Standard `mattn/go-sqlite3` requires CGo (`gcc`) for desktop compilation.
- **Decision:** Selected `modernc.org/sqlite`.
- **Rationale:** 100% pure Go implementation eliminates host C-compiler dependencies on desktop, guarantees lightning-fast compilation, and prevents glibc/musl version mismatches across different Linux distributions. CGo is reserved strictly for Android mobile cross-compilation where the NDK Clang toolchain is explicitly provided.

### ADR-02: Recursive SQLite Triggers for Polymorphic Cascading Deletes
- **Context:** MetaData nodes can be parented by an `Idea`, `Memory`, `InterLink`, or another `MetaData` node (up to 6 levels deep). When an Idea or root node is deleted, all descendant sub-thoughts and attachments must be purged.
- **Decision:** Implemented SQLite trigger `cascade_metadata_delete` with connection string flag `_pragma=recursive_triggers(1)`.
- **Rationale:** Rather than writing complex recursive delete traversals in Go application code (which risk race conditions and partial failures), SQLite natively executes recursive triggers in a single atomic transaction.

### ADR-03: Connection Pool Cursor Deadlock Resolution
- **Context:** In SQLite, iterating `sql.Rows` with `MaxOpenConns=1` holds a lock on the connection. Executing nested queries inside a `rows.Next()` loop caused `database is locked` deadlocks during test runs.
- **Decision:** Modified `GetMetaDataForParent` and `GetInterLinksForSource` to eagerly read all rows into memory slices and close the cursor before executing child queries.
- **Rationale:** Guarantees deterministic query completion in 0.009s with zero connection starvation.

### ADR-04: Google Gemini 2.5 Flash API for AI Expansion
- **Context:** The app required an AI model to assist users in brainstorming sub-thoughts and expanding ideas.
- **Decision:** Integrated `google.golang.org/genai` utilizing the `gemini-2.5-flash` model.
- **Rationale:** Gemini 2.5 Flash provides near-instant response latency with strong structured reasoning capabilities. Users provide their own API key, preserving privacy and keeping IdeaMapper free of ongoing infrastructure fees.

### ADR-05: Git Large File Size Management (100 MB Limit)
- **Context:** `IdeaMapper.apk` is 174 MB (multi-ABI universal APK). Pushing to GitHub caused rejection via pre-receive hook error `remote: error: File IdeaMapper.apk is 173.38 MB; this exceeds GitHub's file size limit of 100.00 MB`.
- **Decision:** Added `.gitignore` to ignore `*.apk`, `*.aab`, and local binaries; removed `IdeaMapper.apk` from git index while keeping the built file intact on disk.
- **Rationale:** Keeps repository clone speeds instant, respects GitHub policies, and allows users to build the APK locally via `./build_apk.sh` or download future releases from GitHub Releases.

### ADR-06: 4-Tab Bottom Navigation with Mobile-First Viewport
- **Context:** Designing an interface that feels equally native on Android touchscreens and desktop monitors.
- **Decision:** Configured Fyne `container.NewAppTabs` with `container.TabLocationBottom` (`Ideas`, `Memories`, `Mind Map`, `Settings`) and initialized desktop window dimensions to `400x700`.
- **Rationale:** Emulates modern mobile UX conventions while scaling cleanly to larger desktop monitors.

### ADR-07: Multi-Layered Living Documentation Enforcement
- **Context:** Project standards require that files in `documents/` (roadmap, setup, prompts/decisions, command history) remain continuously updated across future development turns without manual prompting.
- **Decision:** Implemented a three-tiered enforcement architecture:
  1. **Workspace Agent Rules (`AGENTS.md` & `GEMINI.md`):** Loaded unconditionally into Antigravity's system prompt on every turn, establishing mandatory living documentation rules.
  2. **Lifecycle Hooks (`.agent/hooks.json`):** Configured a `PreInvocation` hook (`scripts/hooks/pre_invocation.sh`) that injects a persistent workspace directive before every turn.
  3. **Automated Logging Utilities (`scripts/log_command.sh`):** Provides a 1-line shell script for atomic, formatted appending to `documents/04_COMMAND_HISTORY.md`.
- **Rationale:** Guarantees that any future AI agent or developer session operating in this repository maintains 100% documentation fidelity and auditability.

