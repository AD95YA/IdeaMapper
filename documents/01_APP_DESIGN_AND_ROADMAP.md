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

## 3. V2 Roadmap (Near-Term Enhancements)

- [ ] **Interactive Visual Graph Canvas**:
  - Implement a 2D node-and-edge visual graph canvas inside the **Mind Map** tab (similar to Obsidian Graph View or Mindly's circular planetary layout).
  - Drag-and-drop nodes, pinch-to-zoom, and visual clustering of connected thoughts.
- [ ] **Automated Google Drive OAuth2 Sync**:
  - Complete the client credential exchange in `api/gdrive.go` to provide automated background cloud snapshots on Wi-Fi.
- [ ] **Full-Text Search (SQLite FTS5)**:
  - Integrate SQLite FTS5 virtual tables to index all Idea descriptions, Memory narratives, and nested MetaData nodes for millisecond fuzzy search.
- [ ] **Rich Markdown Rendering**:
  - Render GitHub-flavored markdown in Idea and Memory detail views (checklists, headers, bullet points, and code syntax highlighting).
- [ ] **Theme Customization**:
  - Support Light, Dark, Solarized, and AMOLED Black themes using Fyne theme overrides.
- [ ] **Biometric / PIN Security**:
  - Optional local PIN or biometric unlock (on Android) to protect private journal entries and sensitive memories.

---

## 4. V3 Roadmap (Long-Term Architectural Vision)

- [ ] **End-to-End Encrypted Peer-to-Peer Sync**:
  - Decentralized local sync between desktop and mobile devices via local network Wi-Fi or WebDAV/libp2p without central servers.
- [ ] **Voice-to-Thought Audio Journaling**:
  - Record voice memos directly inside Memories or Ideas with on-device speech-to-text transcription (via Whisper.cpp).
- [ ] **Local Semantic Search & On-Device Embeddings**:
  - Generate local vector embeddings using small on-device models (e.g., Gemma / Ollama) for semantic concept clustering ("find thoughts related to sustainability").
- [ ] **Universal Export & Interoperability**:
  - Export entire knowledge graphs to standard open formats: Obsidian Markdown Vaults, OPML, GraphML, and printable PDF digests.
- [ ] **Spatial Planetary Navigation (Mindly Classic Paradigm)**:
  - Implement an orbit-based zoom navigation where clicking a parent node centers it as a sun with child sub-thoughts orbiting around it.
