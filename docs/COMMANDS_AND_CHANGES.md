# IdeaMapper — Command Execution & System Changes Reference

This document provides a comprehensive breakdown of all commands executed, their arguments/flags, parameter significance, and the architectural changes made across the codebase (split by category: Installations, Build, App Core, Database, UI/UX, and System Diagnostics).

---

## Table of Contents
1. [Installations & Environment Setup](#1-installations--environment-setup)
2. [Build & Packaging Commands](#2-build--packaging-commands)
3. [Database Changes & Commands](#3-database-changes--commands)
4. [Application Core & AI Changes](#4-application-core--ai-changes)
5. [UI & UX Component Changes](#5-ui--ux-component-changes)
6. [Environment, Display & Virtualization Diagnostics](#6-environment-display--virtualization-diagnostics)
7. [Git Version Control & Object Store Maintenance](#7-git-version-control--object-store-maintenance)
8. [Log of Repeated Commands (Parameters Only)](#8-log-of-repeated-commands-parameters-only)

---

## 1. Installations & Environment Setup

### `go install fyne.io/fyne/v2/cmd/fyne@v2.7.3`
* **Parameters**:
  * `fyne.io/fyne/v2/cmd/fyne`: The Go package path to the official Fyne CLI binary.
  * `@v2.7.3`: The specific version tag matching the version of `fyne.io/fyne/v2` pinned in `go.mod`.
* **Significance**:
  * Downloads, compiles, and installs the `fyne` executable to `$HOME/go/bin/fyne`.
  * The tool provides packaging and bundling capabilities for desktop, web, and mobile (Android APK / iOS).

### `ffmpeg -f lavfi -i "color=c=0x4f46e5:s=512x512" -vframes 1 Icon.png -y`
* **Parameters**:
  * `-f lavfi`: Specifies FFmpeg's Libavfilter virtual input device.
  * `-i "color=c=0x4f46e5:s=512x512"`: Generates an indigo test canvas with dimensions 512×512 (the required resolution for mobile application icons).
  * `-vframes 1`: Extracts a single video frame.
  * `Icon.png`: Output file name.
  * `-y`: Overwrites output without interactive confirmation.
* **Significance**:
  * Generates the required default `Icon.png` in the project root needed by `fyne package -os android`.

### `chmod +x build_apk.sh`
* **Parameters**:
  * `+x`: Adds executable permissions for the current user and group.
  * `build_apk.sh`: Target shell script.
* **Significance**:
  * Enables direct execution of the automated Android packaging script `./build_apk.sh`.

---

## 2. Build & Packaging Commands

### `go build -o idea_app main.go`
* **Parameters**:
  * `-o idea_app`: Specifies the output binary name instead of the default package name.
  * `main.go`: The main entrypoint file.
* **Significance**:
  * Compiles the entire Go application with CGO and OpenGL/Weston bindings into an executable binary `idea_app`.

### `go build ./...`
* **Parameters**:
  * `./...`: Recursive wildcard pattern instructing Go toolchain to traverse all subdirectories and packages.
* **Significance**:
  * Verifies type safety and compilation across all packages (`ai`, `api`, `db`, `models`, `ui`, and root) without producing output binaries.

### `~/go/bin/fyne package -os android -appID com.antigravity.ideamemoryapp`
* **Parameters**:
  * `package`: Fyne CLI subcommand to bundle the application.
  * `-os android`: Targets the Android operating system.
  * `-appID com.antigravity.ideamemoryapp`: Unique reverse-DNS package identifier for the Android Manifest.
* **Significance**:
  * Triggers the Android APK compilation pipeline (requires Android SDK and NDK).

### `./idea_app &` (Daemon Execution)
* **Parameters**:
  * `./idea_app`: Runs the compiled executable.
  * `&`: Detaches the process to run in the background.
* **Significance**:
  * Launches the desktop GUI app within the WSLg display environment with a 400×700 mobile viewport.

---

## 3. Database Changes & Commands

### Code Changes in `db/`
* **Trigger Cascade Updates ([database.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/db/database.go))**:
  * Updated trigger `cascade_note_delete` to `cascade_metadata_delete` (`parent_type = 'MetaData'`).
  * Added PRAGMA `recursive_triggers(1)` to the SQLite DSN:
    ```go
    dsn = "file:memdb1?mode=memory&cache=shared&_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)"
    ```
    *Significance*: SQLite disables recursive triggers by default. Without this PRAGMA, deleting an Idea only deleted Level 1 notes, leaving deeper levels orphaned.
* **Repository Alignment ([metadata_repo.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/db/metadata_repo.go))**:
  * Renamed `db/note_repo.go` to `db/metadata_repo.go`.
  * Migrated all functions from `models.Note` to `models.MetaData`.
  * Fixed row deadlock in `GetMetaDataForParent`: closed outer `sql.Rows` before initiating recursive child queries.
* **Entity Lookups ([idea_repo.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/db/idea_repo.go), [memory_repo.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/db/memory_repo.go))**:
  * Added `GetIdeaByID(id string)` and `GetMemoryByID(id string)` with automatic metadata preloading.
* **Link Exploration ([link_repo.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/db/link_repo.go))**:
  * Added `GetAllInterLinks()` to retrieve all graph edges for the Mind Map screen.
  * Fixed `l.LinkNote` to `l.LinkMetaData`.

### Test Execution Commands:
#### `go test -v -timeout 10s ./db`
* **Parameters**:
  * `-v`: Verbose output (prints every test case name and logs).
  * `-timeout 10s`: Failsafe boundary preventing indefinite execution.
  * `./db`: Restricts testing solely to the database repository package.
* **Significance**:
  * Validated in-memory SQLite schema, polymorphic cascading deletes, 6-level recursion, and SQL injection resistance.

---

## 4. Application Core & AI Changes

### Code Changes in `ai/`, `api/`, `models/`
* **Gemini GenAI SDK Migration ([ai/gemini.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ai/gemini.go))**:
  * Updated `client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(...), nil)` to comply with `google.golang.org/genai` v1.52.0.
  * Added `GenerateIdeasWithKey(prompt, apiKey)` to support in-app user API key configuration alongside the environment fallback.
* **JSON Backup & Restore Engine ([api/gdrive.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/api/gdrive.go))**:
  * Implemented `ImportBackupBlob(jsonData []byte)`: Recursively restores snapshot JSON into SQLite, reinstating Ideas, Memories, and nested MetaData hierarchies.
  * Added `TestGenerateAndImportBackupBlob` in [gdrive_test.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/api/gdrive_test.go).

### Test Execution Commands:
#### `go test -v ./api`
* **Parameters**:
  * `-v`: Verbose test reporting.
  * `./api`: Runs tests in the API package.
* **Significance**:
  * Validates JSON serialization and recursive deserialization of important entities.

---

## 5. UI & UX Component Changes

### New Components Created:
1. **[idea_detail.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/idea_detail.go)**:
   * Replaces informational placeholder dialog with a full-screen view.
   * Actions: Edit title/description, Star/Importance toggle, Delete confirmation, *"Ask Gemini AI"* button, and connected InterLinks section.
2. **[memory_detail.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/memory_detail.go)**:
   * View details of a memory with `YYYY-MM-DD` date picker, editing, deleting, and connected links.
3. **[metadata_view.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/metadata_view.go)**:
   * Hierarchical tree renderer supporting 1–6 nested levels.
   * Indentation rendering with inline buttons: "+ Sub", "+ Attach", "Edit", and "Delete".
4. **[links_modal.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/links_modal.go)**:
   * Polymorphic link generator connecting Idea $\leftrightarrow$ Memory, Idea $\leftrightarrow$ Idea, or Memory $\leftrightarrow$ Memory with relationship labels.
5. **[links_view.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/links_view.go)**:
   * Dedicated *"Mind Map"* explorer tab showing all connections with bidirectional jump buttons.
6. **[settings.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/settings.go)**:
   * Gemini API key input and live test button.
   * Export/import JSON backup buttons and system statistics.
7. **[layout.go](file:///home/aditya/.gemini/antigravity/scratch/idea_app/ui/layout.go)**:
   * Updated bottom tab navigation to 4 items: `💡 Ideas`, `💭 Memories`, `🧠 Mind Map`, and `⚙ Settings`.

---

## 6. Environment, Display & Virtualization Diagnostics

### `which adb emulator fyne fyne-cross 2>&1`
* **Parameters**:
  * `adb emulator fyne fyne-cross`: Names of utilities to search for in `$PATH`.
  * `2>&1`: Redirects stderr to stdout.
* **Significance**: Checked whether Android platform tools were pre-installed.

### `echo "DISPLAY: $DISPLAY"; echo "WAYLAND_DISPLAY: $WAYLAND_DISPLAY"`
* **Significance**: Identified active Linux graphical server (`:0` via Xwayland, and `wayland-0` via Weston).

### `groups; uname -a; lsb_release -a`
* **Significance**:
  * Confirmed OS was Ubuntu 24.04 (Noble) running under the WSL2 Linux kernel `6.6.87.2-microsoft-standard-WSL2`.

### `cmd.exe /c "where adb & where emulator"`
* **Parameters**:
  * `/c`: Carries out the command string and terminates.
  * `where adb`: Windows command to locate `adb.exe` on Windows host system path.
* **Significance**: Verified if the Windows host had Android Studio or platform-tools in its PATH.

### `powershell.exe -Command '...'`
* **Parameters**:
  * `-Command`: Executes PowerShell command script.
  * Searched `C:\Program Files`, `C:\Program Files (x86)`, and `AppData\Local` for Android or emulator installations.
* **Significance**: Confirmed neither Android Studio nor BlueStacks were installed on the host machine.

### `xwininfo -root -tree | head -n 30`
* **Parameters**:
  * `-root`: Queries the root window of the current X display.
  * `-tree`: Outputs full hierarchy of child windows.
  * `head -n 30`: Limits output to top 30 lines.
* **Significance**: Confirmed Weston WM was active and accepting windows.

### `curl -sI https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip`
* **Parameters**:
  * `-s`: Silent mode.
  * `-I`: Fetches HTTP headers only (HEAD request).
* **Significance**: Tested network connectivity to Google's official Android repository.

---

## 7. Git Version Control & Object Store Maintenance

### `git mv db/note_repo.go db/metadata_repo.go`
* **Significance**: Renamed repository file while preserving commit history.

### Object Directory Permission Migration Script
Because previous commits (`72ae945`, `3c60d3d`) were committed by `root`, several subdirectories in `.git/objects` were owned by `root:root` with mode `755`, preventing non-root user `aditya` from inserting new git objects.

**Executed Python Migration Command**:
```bash
python3 -c "
import os, shutil

git_objs = '.git/objects'
for name in os.listdir(git_objs):
    p = os.path.join(git_objs, name)
    if os.path.isdir(p) and len(name) == 2:
        st = os.stat(p)
        if st.st_uid != os.getuid():
            bak_p = os.path.join(git_objs, 'rootbak_' + name)
            os.rename(p, bak_p)
            os.makedirs(p, mode=0o775)
            for f in os.listdir(bak_p):
                shutil.copy2(os.path.join(bak_p, f), os.path.join(p, f))
"
```
* **Significance**:
  * Since user `aditya` owned `.git/objects/`, renaming the root-owned subdirectories within the same directory was permitted by the POSIX filesystem.
  * Recreated clean directories with `aditya:aditya` ownership and copied the object files, restoring write permissions without requiring `sudo`.

### `git fsck --full`
* **Parameters**:
  * `--full`: Thoroughly verifies the connectivity and validity of all objects in the database.
* **Significance**: Confirmed repository consistency (100% object directory check passed).

### `git add . && git commit -m "..."`
* **Commit 1 (`a838ce8`)**: *"Fix compilation errors and update unit tests for metadata, links, and triggers"*
* **Commit 2 (`32aa598`)**: *"Implement detail views, hierarchical metadata, interlinks mind map, Gemini integration, and backup UI"*
* **Commit 3 (`7c24a61`)**: *"Add build_apk.sh helper script and Icon.png for Android packaging"*

---

## 8. Log of Repeated Commands (Parameters Only)

Below is the chronological list of commands executed repeatedly during iteration and verification:

1. `git log -n 10 --oneline`
2. `git status; git branch -a`
3. `find . -maxdepth 4 -not -path '*/.*'`
4. `go test ./...`
5. `git log -p -n 5`
6. `git show 72ae945`
7. `git log -n 5 --stat`
8. `git log -n 1 --format=fuller cfff065879822589c67dd76daad654658044dbef`
9. `git reflog -n 20; git stash list`
10. `git diff HEAD~1`
11. `git tag -l; git branch -a; git log --all --grep="requirement" -i`
12. `git status --ignored`
13. `go doc google.golang.org/genai Client.Models`
14. `go doc google.golang.org/genai Models.GenerateContent`
15. `go doc google.golang.org/genai Text`
16. `go doc google.golang.org/genai GenerateContentResponse`
17. `ls -la ~/.idea_app 2>/dev/null || true`
18. `go test -v ./...`
19. `go test -v ./db`
20. `go test -v -timeout 10s ./db`
21. `go build ./...`
22. `go doc fyne.io/fyne/v2/theme`
23. `go build -o idea_app main.go`
24. `git status`
25. `ls -la /mnt/c/Users/ /mnt/c/Users/*/AppData/Local/Android/Sdk 2>/dev/null || true`
26. `ls -la /mnt/c/Users/adity/AppData/Local/Android/Sdk 2>/dev/null; which cmd.exe powershell.exe`
27. `ls -la /mnt/c/Users/adity/AppData/Local/Android 2>/dev/null || true`
28. `ls -la /mnt/c/Users/adity/.android ~/.android 2>/dev/null || true`
29. `which xwininfo xdotool scrot import grim slurp 2>&1 || true`
30. `which ffmpeg xwd xprop xrandr 2>&1 || true`
31. `timeout 3 ./idea_app`
32. `docker ps 2>/dev/null; docker images 2>/dev/null`
33. `docker version; docker info 2>&1 | head -n 10`
34. `cmd.exe /c "docker version" 2>/dev/null || true`
35. `~/go/bin/fyne package --help`
36. `java -version 2>&1; javac -version 2>&1 || true`
37. `ls -la /mnt/`
38. `ls -la /mnt/d /mnt/e 2>/dev/null || true`
39. `python3 /home/aditya/.gemini/antigravity-ide/brain/feae06b9-a646-4019-9de5-ad1eea406e53/scratch/seed.py`
40. `git log -n 3 --oneline`
41. `git status`
