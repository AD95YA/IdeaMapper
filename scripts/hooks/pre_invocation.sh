#!/usr/bin/env bash
# Antigravity PreInvocation Lifecycle Hook for IdeaMapper
# Consumes JSON on stdin, outputs step injection JSON on stdout.

cat << 'EOF'
{
  "injectSteps": [
    {
      "systemMessage": {
        "systemMessage": "MANDATORY WORKSPACE PROTOCOL: Ensure all files in documents/ (01_APP_DESIGN_AND_ROADMAP.md, 02_SETUP_AND_INSTALLATION.md, 03_PROMPTS_AND_DECISIONS.md, and 04_COMMAND_HISTORY.md) are updated with any actions, decisions, and commands executed during this turn."
      }
    }
  ]
}
EOF
