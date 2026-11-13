#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# IdeaMapper — Command History Logger Helper
# Appends a standardized entry to documents/04_COMMAND_HISTORY.md
#
# Usage:
#   ./scripts/log_command.sh "<command_executed>" "<exit_code>" "<summary_of_changes>"
# Example:
#   ./scripts/log_command.sh "go test -v ./..." "0" "Ran package test suites; all passed in 0.02s."
# ==============================================================================

WORKSPACE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET_FILE="${WORKSPACE_ROOT}/documents/04_COMMAND_HISTORY.md"

if [ $# -lt 3 ]; then
    echo "Usage: $0 <command_executed> <exit_code> <summary_of_changes>"
    exit 1
fi

COMMAND_STR="$1"
EXIT_CODE="$2"
SUMMARY_STR="$3"
TIMESTAMP="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

if [ ! -f "$TARGET_FILE" ]; then
    echo "Error: ${TARGET_FILE} does not exist."
    exit 1
fi

cat << EOF >> "$TARGET_FILE"

### [${TIMESTAMP}] Command Execution
- **Command:** \`${COMMAND_STR}\`
- **Exit Code:** \`${EXIT_CODE}\`
- **Summary:** ${SUMMARY_STR}
EOF

echo "✔ Logged command to documents/04_COMMAND_HISTORY.md: '${COMMAND_STR}' (Exit: ${EXIT_CODE})"
