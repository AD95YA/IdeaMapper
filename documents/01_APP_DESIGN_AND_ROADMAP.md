# IdeaMapper — Application Design, Architecture & Roadmap

## 1. Executive Summary & Vision

**IdeaMapper** (`com.antigravity.ideamemoryapp`) is an offline-first, cross-platform personal knowledge management (PKM), mind mapping, and reflective memory journaling application. Built with **Go** and **Fyne v2**, IdeaMapper bridges the cognitive gap between linear note-taking and spatial graph thinking by treating thoughts, memories, and their multi-tier sub-components as interconnected nodes in a personal knowledge graph.

### The Problem It Solves
Human brains think associatively and hierarchically. Standard note-taking software forces users into rigid flat lists or isolated text documents. Meanwhile, commercial mind-mapping applications (most notoriously **Mindly** on Android) rug-pulled paying users by removing lifetime purchases, killing original applications, and introducing stripped-down subscription paywalls. IdeaMapper was built as an enduring, anti-greed antidote: 100% offline-first, permanent, free of subscriptions or tracking, and entirely controlled by the user.

---

## 2. V1 MVP Scope (Fully Implemented)

### Core Modules & Deliverables
1. **💡 Ideas Hub (`ui/ideas_list.go`, `ui/idea_detail.go`)**:
   - Rapid capture of ideas with titles, detailed text descriptions, icons, and importance flags.
   - Real-time search by keyword across titles and bodies.
   - Instant filtering for **★ Starred Only** high-priority initiatives.
   - Comprehensive detail view with in-place editing, deletion confirmation, and embedded action triggers.

2. **💭 Memories & Timeline Journal (`ui/memories_list.go`, `ui/memory_detail.go`)**:
   - Date-anchored personal memories organized chronologically and grouped by year.
   - Starred status for significant life milestones.
   - Bi-directional association linking memories to the future ideas they sparked.

3. **🌱 Recursive Multi-Level Sub-Thoughts (`ui/metadata_view.go`, `db/metadata_repo.go`)**:
   - Polymorphic tree structure supporting up to 6 levels of nesting (`Level` 1..6).
   - Parent entities can be an `Idea`, `Memory`, `InterLink`, or another `MetaData` node.
   - Multi-attachment support for file paths (`📎`) and photos (`📷`).
   - SQLite recursive triggers guaranteeing clean cascading deletions down through all sub-levels.

4. **🧠 InterLinks Mind Map (`ui/links_view.go`, `ui/links_modal.go`, `db/link_repo.go`)**:
   - Typed, bidirectional associative links connecting heterogeneous entities (`Idea ↔ Memory`, `Idea ↔ Idea`, `Memory ↔ Memory`).
   - Custom relationship tagging (e.g., *"Inspired By"*, *"Depends On"*, *"Learned From"*).
   - Embedded link metadata for attaching specific notes to the relationship itself.
   - Dedicated connection search explorer with one-tap deep navigation into source and target nodes.

5. **🤖 Google Gemini 2.5 Flash AI (`ai/gemini.go`)**:
   - One-tap creative idea expansion directly inside Idea detail views.
   - Prompts Gemini 2.5 Flash with current context to generate structured sub-thoughts, questions, and implementation avenues.
   - In-app API key configuration with live connection testing and zero telemetry.

6. **⚙️ Data Sovereignty & Portability (`ui/settings.go`, `api/gdrive.go`)**:
   - Offline-first SQLite database (`~/.idea_app/app.db` or Android internal storage).
   - Portable JSON snapshot export and import (`BackupPayload`) containing ideas, memories, and multi-tier sub-metadata.
   - Cloud synchronization scaffolding for Google Drive OAuth integration.
   - Real-time system statistics (Total Ideas, Memories, and InterLinks).

7. **📱 Android Mobile & Desktop Parity (`build_apk.sh`, `main.go`)**:
   - Desktop application running on Linux (X11/Wayland), macOS, and Windows with a responsive 400×700 mobile viewport.
   - Full Android APK packaging pipeline producing standalone `IdeaMapper.apk` ready for sideloading via ADB.

---

## 3. V2 Scope (Delivered & Verified)

- [x] **Mobile Navigation & Breadcrumbs (`ui/nav_stack.go`)**:
  - Unified persistent `<` (Back) and `⌂` (Home) controls across all views.
  - Android hardware/gesture back key interception (`KeyEscape`) mapping to stack traversal.
  - Horizontally scrollable breadcrumb trail (`container.NewHScroll`) supporting deep hierarchies up to 6 levels without clipping.
- [x] **Swipeable Tabs & Conflict Guard (`ui/swipe_container.go`)**:
  - Horizontal drag gesture switching across all 4 tabs (`Ideas` ↔ `Memories` ↔ `Mind Map` ↔ `Settings`).
  - Gesture conflict guard on the Mind Map tab reserving center canvas drag for panning while allowing edge-swipes (`< 24dp`) and header swipes for tab navigation.
- [x] **Performance & Memory Optimization (`ui/widget_pool.go`, `ui/image_cache.go`, `ui/metadata_view.go`)**:
  - Active-orbit culling: only focused nodes and direct children are laid out; deep descendant sub-trees lazy-load on demand.
  - Zero-allocation graphic primitive pooling (`CirclePool`, `LinePool`, `TextPool`).
  - Async thumbnail downscaling (128x128 for orbit nodes, 512x512 for detail views) with a bounded 50-item LRU cache.
- [x] **Upgraded Mind Map UI (`ui/mindmap_ascii.go`, `ui/mindmap_radial.go`, `ui/links_view.go`)**:
  - Rich Unicode box-drawing tree view with status badges: `[★]` Starred, `[📝]` Notes, `[🖼️]` Attachments, `[⏳]` Solo/Unlinked thought.
  - Radial orbit overlay with directional connector lines and depth-graded accent colors (Levels 1–6).
  - 3-way top view switcher: Box-Drawing Tree, Radial Orbit, and Connection Cards.
- [x] **Automated Daily Local & SQLite Backups (`api/backup_scheduler.go`, `ui/settings.go`)**:
  - Startup and resume scheduler automatically creating daily JSON snapshots and raw SQLite `.db` clones.
  - Rolling 7-day retention policy auto-pruning older archives.
  - Manual "Backup Now" trigger and last backup status in Settings.
- [x] **First-Launch JSON Onboarding Restore (`ui/onboarding.go`)**:
  - Clean onboarding dialog prompting first-time users to restore an existing JSON backup or start fresh with a starter topic.
- [x] **Google Sign-In & Drive Sync Preview (`api/google_auth.go`, `ui/settings.go`)**:
  - Google OAuth profile caching and signed-in status display with Sign Out option.
  - "Google Drive Cloud Sync — Coming Soon" informational modal.
- [x] **Gemini API Onboarding & Copyable Prompt (`ui/settings.go`)**:
  - Expandable 3-step guide for getting a free Gemini API key from Google AI Studio.
  - Read-only prompt template preview with one-tap "Copy Prompt" clipboard action.
- [x] **Comprehensive Test Suite & V2 Release (`releases/ideamapper-v2.apk`)**:
  - Complete edge-case unit test coverage passing with `-race` across all packages.
  - Release-ready universal APK package built via `./scripts/build_and_release_apk.sh v2`.

---

## 4. V3 Roadmap (Future Architectural Vision)

- [ ] **Automated Google Drive OAuth2 Sync**:
  - Complete background cloud upload via Drive AppData API.
- [ ] **Full-Text Search (SQLite FTS5)**:
  - Virtual tables indexing all descriptions, notes, and attachment metadata.
- [ ] **Rich Markdown Rendering**:
  - GitHub-flavored markdown with checklists and code blocks.
- [ ] **Theme Customization**:
  - Dark, Light, Solarized, and AMOLED Black themes.
- [ ] **Biometric / PIN Security**:
  - Local biometric unlock on Android for sensitive memories.
- [ ] **End-to-End Encrypted Peer-to-Peer Sync**:
  - Decentralized local sync between desktop and mobile devices via local network Wi-Fi or WebDAV/libp2p without central servers.
- [ ] **Voice-to-Thought Audio Journaling**:
  - Record voice memos directly inside Memories or Ideas with on-device speech-to-text transcription (via Whisper.cpp).
- [ ] **Local Semantic Search & On-Device Embeddings**:
  - Generate local vector embeddings using small on-device models (e.g., Gemma / Ollama) for semantic concept clustering ("find thoughts related to sustainability").
- [ ] **Universal Export & Interoperability**:
  - Export entire knowledge graphs to standard open formats: Obsidian Markdown Vaults, OPML, GraphML, and printable PDF digests.

