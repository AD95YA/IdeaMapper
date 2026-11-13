# IdeaMapper — Project Rules & System Architecture

- **Title & Package:** `IdeaMapper` (Short title: `IdeaMapper`, Android package / Application ID: `com.antigravity.ideamemoryapp`).
- **Tech Stack:** Go (Golang 1.22+) + Fyne v2 (`fyne.io/fyne/v2`) targeting Desktop (Linux X11/Wayland, Windows, macOS) and Android (`fyne package -os android`).
- **Database & Storage:** Pure Go embedded SQLite via `modernc.org/sqlite` (100% CGo-free). Foreign keys (`_pragma=foreign_keys(1)`) and recursive triggers (`_pragma=recursive_triggers(1)`) are strictly enforced. Storage directory is resolved via Fyne URI storage root on Android and fallback to `~/.idea_app/app.db` on Desktop.
- **UI & Layout Design:** 4-tab bottom navigation (`💡 Ideas`, `💭 Memories`, `🧠 Mind Map`, `⚙ Settings`). Viewport is initialized to `400x700` logical resolution on desktop for seamless mobile-to-desktop parity, with responsive resizing across larger screens.
- **Data Hierarchy & Graph Topology:**
  - **Polymorphic MetaData Tree:** Hierarchical sub-thoughts supporting up to 6 levels of nesting (`Level` 1..6). Parent entities can be an `Idea`, `Memory`, `InterLink`, or another `MetaData` node. Nodes support external file paths (`📎`) and photos (`📷`).
  - **InterLinks Associative Graph:** Typed, bidirectional connections between heterogeneous entities (`Idea ↔ Memory`, `Idea ↔ Idea`, `Memory ↔ Memory`) with custom relationship names and embedded link metadata.
  - **Cascading Integrity:** Automated cascade deletes via SQLite triggers (`cascade_idea_delete`, `cascade_memory_delete`, `cascade_interlink_delete`, `cascade_metadata_delete`) guarantee zero orphaned child records or attachments.
- **AI Integration:** Google Gemini 2.5 Flash API (`google.golang.org/genai`) for context-aware idea ideation and creative sub-thought expansion. Dynamic API key configuration via in-app preferences or `GEMINI_API_KEY` environment variable.
- **Data Sovereignty & Portability:** 100% offline-first architecture with zero tracking, local JSON snapshot export/import (`BackupPayload` preserving full nested metadata trees), and Google Drive OAuth cloud backup scaffolding.
- **Anti-Exploit / Anti-Greed Ethos:** No subscriptions, no ads, no telemetry, no artificial feature gates, and strict user ownership of personal knowledge databases.

# Constants
- **AppTitle:** IdeaMapper
- **ShortTitle:** IdeaMapper
- **TagLine:** Map your thoughts, preserve your memories, and connect the sparks of your mind.
- **PackageName:** "com.antigravity.ideamemoryapp"
- **DBFileName:** "app.db"
- **DefaultWindowWidth:** 400
- **DefaultWindowHeight:** 700
- **MaxMetaDataLevel:** 6
- **GeminiModel:** "gemini-2.5-flash"
- **BackupFileName:** "ideamapper_backup.json"
- **TargetAndroidAPI:** 34
- **RecommendedNDK:** "26.3.11579264" (r26d)

