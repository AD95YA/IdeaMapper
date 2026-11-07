#!/usr/bin/env bash
set -e

echo "============================================="
echo "   IdeaMapper - Android APK Build Tool       "
echo "============================================="

# 1. Detect Android SDK
if [ -z "$ANDROID_HOME" ]; then
    # Check common Linux paths
    if [ -d "$HOME/Android/Sdk" ]; then
        export ANDROID_HOME="$HOME/Android/Sdk"
    # Check Windows Host paths via WSL2
    elif [ -d "/mnt/c/Users/$USER/AppData/Local/Android/Sdk" ]; then
        export ANDROID_HOME="/mnt/c/Users/$USER/AppData/Local/Android/Sdk"
    elif compgen -G "/mnt/c/Users/*/AppData/Local/Android/Sdk" > /dev/null; then
        export ANDROID_HOME=$(ls -d /mnt/c/Users/*/AppData/Local/Android/Sdk | head -n 1)
    fi
fi

if [ -n "$ANDROID_HOME" ]; then
    echo "✔ Found Android SDK: $ANDROID_HOME"
    export PATH="$PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools"
else
    echo "⚠ ANDROID_HOME is not set."
fi

# 2. Detect Android NDK
if [ -z "$ANDROID_NDK_HOME" ] && [ -n "$ANDROID_HOME" ]; then
    if [ -d "$ANDROID_HOME/ndk-bundle" ]; then
        export ANDROID_NDK_HOME="$ANDROID_HOME/ndk-bundle"
    elif compgen -G "$ANDROID_HOME/ndk/*" > /dev/null; then
        export ANDROID_NDK_HOME=$(ls -d $ANDROID_HOME/ndk/* | sort -V | tail -n 1)
    fi
fi

if [ -n "$ANDROID_NDK_HOME" ]; then
    echo "✔ Found Android NDK: $ANDROID_NDK_HOME"
else
    echo "⚠ ANDROID_NDK_HOME is not set."
fi

# 3. Ensure Fyne CLI is in PATH
if ! command -v fyne &> /dev/null; then
    if [ -f "$HOME/go/bin/fyne" ]; then
        export PATH="$PATH:$HOME/go/bin"
    else
        echo "Installing Fyne CLI..."
        go install fyne.io/fyne/v2/cmd/fyne@v2.7.3
        export PATH="$PATH:$HOME/go/bin"
    fi
fi
echo "✔ Fyne CLI ready: $(which fyne)"

# 4. Check requirements before build
if [ -z "$ANDROID_NDK_HOME" ]; then
    echo ""
    echo "------------------------------------------------------------------------"
    echo "To build the APK file for a physical device, Android NDK is required."
    echo ""
    echo "Step 1: If using Android Studio:"
    echo "  Open Android Studio -> Tools -> SDK Manager -> SDK Tools tab"
    echo "  Check 'NDK (Side by side)' and 'Android SDK Command-line Tools', click Apply."
    echo ""
    echo "Step 2: Export paths or re-run this script:"
    echo "  export ANDROID_HOME=\"/path/to/Sdk\""
    echo "  export ANDROID_NDK_HOME=\"\$ANDROID_HOME/ndk/<version>\""
    echo "  ./build_apk.sh"
    echo "------------------------------------------------------------------------"
    exit 1
fi

echo "Building IdeaMapper Android APK..."
fyne package -os android -appID com.antigravity.ideamemoryapp -name "IdeaMapper" -icon Icon.png

echo ""
echo "✔ Build complete! APK generated: IdeaMapper.apk"
echo "You can install it to your connected device using: adb install IdeaMapper.apk"
