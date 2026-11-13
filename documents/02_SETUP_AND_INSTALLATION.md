# IdeaMapper — Setup, Toolchain & Installation Guide

## 1. Development Environment Overview

IdeaMapper is engineered and verified on an **Ubuntu 24.04 LTS (Noble Numbat)** environment running on **WSL2** under Windows 11, with cross-platform targets for Desktop (Linux/Windows/macOS) and Android.

### Installed Tooling & Verified Versions
| Tool / Component | Version | Installation Path / Identifier |
| :--- | :--- | :--- |
| **Operating System** | Ubuntu 24.04.4 LTS (WSL2) | Linux Kernel 5.15+ |
| **Go (Golang)** | `go version go1.22.2 linux/amd64` | `/usr/bin/go` |
| **Fyne GUI Toolkit** | `v2.7.3` | `fyne.io/fyne/v2` |
| **Fyne CLI** | `v2.7.3` | `$HOME/go/bin/fyne` |
| **Java Development Kit** | OpenJDK `21.0.12.1` | `/usr/bin/java` |
| **Android SDK Tools** | `sdkmanager` v12.0 | `$HOME/Android/Sdk/cmdline-tools/latest/bin/sdkmanager` |
| **Android Platform** | Android 14 (`platforms;android-34`) | `$HOME/Android/Sdk/platforms/android-34` |
| **Android Build-Tools**| `34.0.0` (`aapt2`, `zipalign`) | `$HOME/Android/Sdk/build-tools/34.0.0` |
| **Android NDK** | `26.3.11579264` (NDK r26d LTS) | `$HOME/Android/Sdk/ndk/26.3.11579264` |
| **Android Platform-Tools** | `adb`, `fastboot` | `$HOME/Android/Sdk/platform-tools` |
| **Display Server** | Xwayland (`DISPLAY=:0`), Wayland (`wayland-0`) | WSLg Weston Desktop Integration |

---

## 2. Desktop Prerequisites & Dependencies

### Ubuntu / Debian Linux Host
To compile Fyne applications with native OpenGL/X11 bindings on Linux:

```bash
sudo apt update
sudo apt install -y libgl1-mesa-dev xorg-dev gcc pkg-config
```

### Clone & Download Dependencies
```bash
git clone https://github.com/abhokarikar/IdeaMapper.git
cd IdeaMapper
go mod download
```

---

## 3. Desktop Build & Execution Commands

### Run Directly
```bash
go run main.go
```

### Build Standalone Executable
```bash
go build -o idea_app main.go
./idea_app
```
*Note: The desktop window initializes at `400x700` logical resolution to emulate a clean mobile screen aspect ratio, but can be freely resized.*

### Run Complete Unit Test Suite
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

---

## 4. Android SDK & NDK Toolchain Installation

Android compilation requires the Android SDK, platform tools, build tools, and the LLVM Clang-based NDK toolchain for CGo cross-compilation.

### Step 1: Download Android Command-Line Tools
```bash
mkdir -p $HOME/Android/Sdk/cmdline-tools
curl -L -o /tmp/cmdline-tools.zip https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip
unzip -q -o /tmp/cmdline-tools.zip -d /tmp/cmdline-extracted
mv /tmp/cmdline-extracted/cmdline-tools $HOME/Android/Sdk/cmdline-tools/latest
rm -rf /tmp/cmdline-extracted /tmp/cmdline-tools.zip
```

### Step 2: Install Required SDK Components & NDK r26d
```bash
yes | $HOME/Android/Sdk/cmdline-tools/latest/bin/sdkmanager --licenses
$HOME/Android/Sdk/cmdline-tools/latest/bin/sdkmanager \
  "platform-tools" \
  "platforms;android-34" \
  "build-tools;34.0.0" \
  "ndk;26.3.11579264"
```

### Step 3: Install Fyne CLI
```bash
go install fyne.io/fyne/v2/cmd/fyne@v2.7.3
```

### Step 4: Configure Persistent Environment Paths
Add the following to your `~/.bashrc`:
```bash
# Android SDK & NDK Configuration
export ANDROID_HOME="$HOME/Android/Sdk"
export ANDROID_NDK_HOME="$ANDROID_HOME/ndk/26.3.11579264"
export PATH="$PATH:$HOME/go/bin:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools:$ANDROID_HOME/build-tools/34.0.0"
```
Reload environment:
```bash
source ~/.bashrc
```

---

## 5. Android APK Build & Sideloading

IdeaMapper provides an automated build tool ([`build_apk.sh`](file:///home/aditya/ai_projects/idea_app/build_apk.sh)) that discovers SDK/NDK paths and triggers Fyne packaging:

```bash
chmod +x build_apk.sh
./build_apk.sh
```

### Manual Packaging Command
```bash
fyne package -os android -appID com.antigravity.ideamemoryapp -name "IdeaMapper" -icon Icon.png
```

### Sideload to Device via ADB
Connect your physical Android phone (with USB Debugging enabled) or start an emulator:
```bash
# Verify connected device
adb devices

# Install APK
adb install -r IdeaMapper.apk
```

---

## 6. Configuration & Runtime Storage

- **Desktop Database Location:** `~/.idea_app/app.db`
- **Android Database Location:** Resolves via internal app storage sandbox URI (`/data/user/0/com.antigravity.ideamemoryapp/files/app.db`).
- **Gemini AI API Key:** Configurable in-app under **Settings ➔ Gemini AI Configuration** or via environment:
  ```bash
  export GEMINI_API_KEY="your-gemini-api-key"
  ```
