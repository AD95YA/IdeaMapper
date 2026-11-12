# IdeaMapper 💡🧠💭

[![Developed with Antigravity IDE](https://img.shields.io/badge/Developed%20with-Antigravity%20IDE-4285F4?style=for-the-badge&logo=google&logoColor=white)](https://deepmind.google)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Fyne GUI](https://img.shields.io/badge/GUI-Fyne%20v2-29BEB0?style=for-the-badge)](https://fyne.io)
[![SQLite](https://img.shields.io/badge/Database-SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)](https://sqlite.org)
[![Gemini AI](https://img.shields.io/badge/AI-Google%20Gemini%202.5%20Flash-8E75B2?style=for-the-badge&logo=googlebard&logoColor=white)](https://deepmind.google/technologies/gemini/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Android%20%7C%20Windows%20%7C%20macOS-informational?style=for-the-badge)](#)

> **Map your thoughts, preserve your memories, and connect the sparks of your mind.**  
> IdeaMapper is an offline-first personal knowledge management, mind mapping, and memory journaling application built with **Go** and **Fyne**, enhanced with **Google Gemini 2.5 Flash AI**, and developed end-to-end using the **Antigravity IDE**.

---

## 📖 Table of Contents

- [Overview & Purpose](#-overview--purpose)
- [Why IdeaMapper? (How Useful It Can Be)](#-why-ideamapper-how-useful-it-can-be)
- [Key Features](#-key-features)
- [How to Best Use IdeaMapper (Workflow Guide)](#-how-to-best-use-ideamapper-workflow-guide)
- [Architecture & Tech Stack](#-architecture--tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Run on Desktop](#run-on-desktop-linux--macos--windows)
  - [Run Unit Tests](#run-unit-tests)
  - [Build Android APK](#build-android-apk)
- [Configuration & Settings](#-configuration--settings)
- [Developed with Antigravity IDE](#-developed-with-antigravity-ide)
- [License](#-license)

---

## 🎯 Overview & Purpose

Human cognition is not a flat list or an isolated set of text files. When we brainstorm or reflect, our thoughts branch into deep sub-concepts, link back to past experiences, and spawn new questions. Traditional note-taking tools force users into rigid hierarchies or endless disconnected lists.

**IdeaMapper** (`com.antigravity.ideamemoryapp`) bridges this divide by providing a unified, graph-aware space to:
1. **Capture raw ideas** whenever inspiration strikes.
2. **Chronicle life memories and events** anchored across time.
3. **Deconstruct thoughts hierarchically** up to 6 levels of recursive sub-metadata with attachments.
4. **Interlink heterogeneous concepts** (Idea ↔ Memory, Idea ↔ Idea, Memory ↔ Memory) into an associative mind map.
5. **Amplify creativity with AI** by using Google Gemini 2.5 Flash directly inside any thought to spark fresh angles and actionable breakdowns.

---

## 🌟 Why IdeaMapper? (How Useful It Can Be)

| Use Case | How IdeaMapper Helps |
| :--- | :--- |
| **Creative Thinkers & Writers** | Structure chaotic brain dumps into hierarchical layers of sub-thoughts, attach reference photos or documents, and let Gemini AI brainstorm unexplored plot hooks or thesis angles. |
| **Engineers & Architects** | Design systems from high-level concepts down to multi-tier technical specifications with attached blueprints and file paths. |
| **Journalers & Reflective Thinkers** | Maintain a chronological timeline of life memories and link them directly to the goals and ideas they inspired. |
| **Lifelong Learners & Researchers** | Map cross-disciplinary connections across different subjects, preserving context through named interlinks. |
| **Privacy & Sovereignty Advocates** | Enjoy 100% offline-first local SQLite storage with zero telemetry, complete with single-click JSON snapshot exports and optional Google Drive synchronization. |

---

## ✨ Key Features

### 1. 💡 Ideas Hub
- **Rapid Capture**: Create new ideas with titles, detailed markdown/plain text, icons, and priority flags.
- **Instant Search & Filtering**: Live search by keywords across titles and body descriptions, with a one-click **★ Starred Only** filter for high-priority initiatives.
- **Interactive Idea Cards**: View summaries and jump straight into full editing or exploration modes.

### 2. 💭 Memories & Timeline Journal
- **Chronological Grouping**: View life events organized by year and date in a readable timeline.
- **Emotional & Historical Context**: Star significant milestones, record detailed narratives, and attach related artifacts.
- **Bi-directional Associative Recall**: See every project or idea that stemmed from a specific memory.

### 3. 🌱 Recursive Multi-Level Sub-Thoughts (MetaData Tree)
- **Deep Nesting (Levels 1 to 6)**: Break any Idea, Memory, or Sub-Thought into child notes without arbitrary depth constraints.
- **Attachment Support**: Associate external file paths (`📎`) and photos (`📷`) with any level of the hierarchy.
- **Safe Cascading Deletions**: Powered by SQLite recursive triggers, ensuring clean cleanup of child branches when a parent is removed.

### 4. 🧠 InterLinks Mind Map Explorer
- **Heterogeneous Connections**: Connect an Idea to a Memory, an Idea to another Idea, or a Memory to a Memory.
- **Named Relationships**: Assign descriptive relationship names (e.g., *"Inspired By"*, *"Depends On"*, *"Contradicts"*, *"Sub-component of"*).
- **Embedded Relationship Notes**: Each link can carry its own dedicated metadata.
- **Graph Search**: Search connections by participant names or relationship types, with one-tap deep navigation into either connected node.

### 5. 🤖 Google Gemini 2.5 Flash AI Integration
- **Context-Aware Ideation**: Tap **"Ask Gemini AI"** on any Idea to automatically analyze its context and generate creative avenues, sub-thoughts, or solutions.
- **Configurable API Key**: Securely input and test your Gemini API key via the in-app Settings panel or pass it via the `GEMINI_API_KEY` environment variable.

### 6. ⚙️ Data Privacy, Backup & Google Drive Sync
- **Local JSON Export**: Generate structured JSON snapshots containing all starred ideas, memories, and multi-tier sub-metadata.
- **One-Click JSON Import**: Restore or merge existing JSON backups seamlessly into the SQLite database.
- **Google Drive Ready**: Built-in cloud sync scaffolding to store backup payloads securely in personal Google Drive storage.
- **Live Diagnostics**: Built-in system counters for total ideas, memories, and active links.

---

## 🚀 How to Best Use IdeaMapper (Workflow Guide)

To get the most out of IdeaMapper, follow this proven 5-step knowledge workflow:

```
[ Step 1: Capture ] ──> [ Step 2: Deconstruct ] ──> [ Step 3: Connect ] ──> [ Step 4: Expand ] ──> [ Step 5: Preserve ]
  Raw Ideas &             Recursive Sub-Notes         InterLink Mind Map        Gemini AI Prompts        Local/Cloud JSON
  Timeline Memories       Levels 1 through 6          Idea ⟷ Memory             Creative Sparks          Backup Snapshots
```

### Step 1: Capture Without Friction
- Whenever an idea or memory strikes, tap **Create Idea** (💡) or **Record Memory** (💭).
- Don't worry about perfection; capture the essence in the title and a brief description.
- Mark critical entries with the **★ Star** so they stand out in filtered views.

### Step 2: Deconstruct Complex Thoughts
- Open any Idea or Memory detail view.
- Under **Sub-Thoughts & Metadata**, click **Add Sub-Thought / Note**.
- Break the thought down into components (e.g., *Level 1: Market Research* ➔ *Level 2: Competitor Analysis* ➔ *Level 3: Feature Gaps*).
- Reference local design files or photos to keep everything in one context.

### Step 3: Weave the Mental Web (InterLinks)
- Navigate to the **InterLinks** section or the **Mind Map** tab (🧠).
- Connect current ideas to past memories that inspired them, or cross-link related concepts.
- Give the link an explicit relationship name (e.g., *"Feasibility Study For"* or *"Learned From Past Failure"*).

### Step 4: Break Mental Blocks with Gemini AI
- When stuck on an idea, open the detail view and click **🤖 Ask Gemini AI**.
- The app sends the title and description to Gemini 2.5 Flash, returning fresh perspectives, potential pitfalls, and creative extensions.
- Save the most promising suggestions directly into your sub-thought hierarchy.

### Step 5: Safeguard & Sync Your Data
- Periodically navigate to **Settings** (⚙).
- Click **Export JSON Backup** to save an offline snapshot of your important ideas and memories.
- Transfer this file across devices or import it anytime with **Import JSON Backup**.

---

## 🛠 Architecture & Tech Stack

```
idea_app/
├── main.go                     # Entrypoint & desktop/mobile storage resolution
├── models/                     # Core domain entities (Idea, Memory, InterLink, MetaData)
├── db/                         # ModernC SQLite driver, polymorphic schema & recursive triggers
├── ui/                         # Cross-platform Fyne v2 GUI tabs, modals, and views
├── ai/                         # Google Gemini 2.5 Flash API client
├── api/                        # Backup/restore serializers & Google Drive uploaders
└── build_apk.sh                # Automated Android packaging script
```

- **Runtime**: [Go (Golang)](https://golang.org) (v1.22+)
- **Cross-Platform GUI**: [Fyne v2](https://fyne.io) (OpenGL / Vulkan mobile & desktop canvas)
- **Database Engine**: [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) (Pure Go, 100% CGo-free SQLite engine)
- **AI Synthesis**: [google.golang.org/genai](https://pkg.go.dev/google.golang.org/genai) (`gemini-2.5-flash`)
- **Android Target**: Android API 34 (Android 14) with NDK r26d Clang toolchain
- **Desktop Targets**: Linux (X11 / Wayland), macOS, and Windows

---

## 📁 Project Structure

```bash
idea_app/
├── Icon.png                 # 512x512 app branding icon
├── IdeaMapper.apk           # Generated Android APK (built via ./build_apk.sh)
├── build_apk.sh             # Interactive / automated Android NDK & SDK packager
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── main.go                  # Application lifecycle and window orchestration
├── ai/
│   └── gemini.go            # Gemini 2.5 Flash client and prompt generation logic
├── api/
│   ├── gdrive.go            # Backup serialization, deserialization & cloud sync
│   └── gdrive_test.go       # Backup payload round-trip unit tests
├── db/
│   ├── database.go          # SQLite schema, indexes, and recursive delete triggers
│   ├── idea_repo.go         # Idea CRUD operations
│   ├── memory_repo.go       # Memory CRUD operations
│   ├── link_repo.go         # InterLink CRUD operations
│   ├── metadata_repo.go     # Recursive multi-level MetaData operations
│   └── repository_test.go   # DB transactions, cascades, and recursive trigger tests
├── docs/
│   └── COMMANDS_AND_CHANGES.md  # Comprehensive chronological architecture log
├── models/
│   ├── models.go            # Domain structs (Idea, Memory, InterLink, MetaData)
│   └── uuid.go              # Crypto-grade random UUID generator
└── ui/
    ├── layout.go            # 4-tab bottom navigation layout
    ├── ideas_list.go        # Ideas dashboard, cards, and search
    ├── idea_detail.go       # Idea details, action bar & AI modal launcher
    ├── memories_list.go     # Memories timeline view with year-based grouping
    ├── memory_detail.go     # Memory detail and editing modal
    ├── metadata_view.go     # Recursive tree view for multi-level sub-thoughts
    ├── links_view.go        # Mind Map global connection browser
    ├── links_modal.go       # Link creation and navigation dialogs
    └── settings.go          # Gemini API key config, JSON backup & Drive sync
```

---

## 💻 Getting Started

### Prerequisites

Ensure you have [Go](https://go.dev/doc/install) (version 1.22 or later) installed on your system.

On Linux (Ubuntu/Debian), ensure graphics development libraries are available for Fyne:
```bash
sudo apt update
sudo apt install -y libgl1-mesa-dev xorg-dev gcc
```

### Run on Desktop (Linux / macOS / Windows)

Clone the repository and launch the app in desktop testing mode:

```bash
# Navigate to workspace
cd idea_app

# Fetch dependencies
go mod download

# Run directly
go run main.go
```

To build a standalone desktop executable:
```bash
go build -o idea_app main.go
./idea_app
```

> **Tip**: On desktop, the window initializes in a convenient 400×700 mobile viewport for responsive testing, but can be resized freely.

### Run Unit Tests

Execute the comprehensive test suite verifying the database triggers, cascade deletes, and JSON backup integrity:

```bash
go test -v ./...
```

Expected output:
```
=== RUN   TestGenerateAndImportBackupBlob
--- PASS: TestGenerateAndImportBackupBlob (0.01s)
=== RUN   TestMainDB
--- PASS: TestMainDB (0.01s)
PASS
ok      idea_app/api
ok      idea_app/db
```

### Build Android APK

IdeaMapper includes an automated build script (`build_apk.sh`) that checks for the Android SDK, Android NDK, and Fyne CLI:

1. **Ensure Android SDK & NDK are installed** (API 34 + NDK r26d recommended):
   ```bash
   export ANDROID_HOME="$HOME/Android/Sdk"
   export ANDROID_NDK_HOME="$ANDROID_HOME/ndk/26.3.11579264"
   export PATH="$PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools"
   ```

2. **Run the build script**:
   ```bash
   ./build_apk.sh
   ```

3. **Install to your Android phone / tablet**:
   ```bash
   adb install -r IdeaMapper.apk
   ```

---

## ⚙️ Configuration & Settings

| Setting | Configuration Method | Description |
| :--- | :--- | :--- |
| **Gemini API Key** | In-App **Settings ➔ Gemini AI Configuration** or `export GEMINI_API_KEY="your-key"` | Unlocks AI-generated sub-thoughts via Gemini 2.5 Flash. |
| **Database Storage** | Automatic | Saved at `~/.idea_app/app.db` on Linux/macOS/Windows, or in the platform sandbox on Android (`com.antigravity.ideamemoryapp`). |
| **JSON Backup Path** | In-App File Chooser (Default: `~/ideamapper_backup.json`) | Target destination for offline portable snapshots. |

---

## ⚡ Developed with Antigravity IDE

IdeaMapper was designed, developed, architected, and packaged using **Antigravity IDE** — Google DeepMind's advanced agentic engineering platform.

### How Antigravity Accelerated IdeaMapper:
- **Autonomous Architecture & Refactoring**: Solved complex polymorphic schema migrations (converting flat notes into recursive 6-level sub-metadata with SQLite triggers) while preserving database integrity and zeroing out compile failures.
- **Full Cross-Package Test Automation**: Built unit test suites validating SQLite memory DSN configurations, foreign key cascade mechanics, and JSON import/export schemas.
- **Native GUI Engineering with Fyne**: Designed a responsive, mobile-first interface featuring 4 primary bottom navigation tabs, custom search filters, dialog workflows, and modal navigation.
- **End-to-End Android Compilation Pipeline**: Diagnosed cross-compilation requirements, installed Android SDK tools and NDK r26d Clang toolchains, configured environment paths, and packaged the production `IdeaMapper.apk`.
- **Chronological Accountability**: Maintained a complete, verified command audit log in `docs/COMMANDS_AND_CHANGES.md`.

---

## 📄 License

This project is licensed under the MIT License — see your repository license configuration for details.

---

<p align="center">
  <b>Built with ❤️ using Go, Fyne, and Antigravity IDE.</b><br>
  <i>Never lose a great idea or a treasured memory again.</i>
</p>
