# Troubleshooting

> 📖 Back to [README](../README.md)

Common issues and fixes for Gentleman Guardian Angel.

---

## "Provider not found"

```bash
# Check if your provider CLI is installed and in PATH
which claude   # Should show: /usr/local/bin/claude or similar
which gemini
which codex
which ollama

# Test if the provider works
echo "Say hello" | claude --print

# For LM Studio, check if the API is accessible
curl http://localhost:1234/v1/models
```

---

## "Rules file not found"

The tool requires a rules file to know what to check:

```bash
# Create your rules file
touch AGENTS.md

# Add your coding standards
echo "# My Coding Standards" > AGENTS.md
echo "- No console.log in production" >> AGENTS.md
```

---

## "Ambiguous response" in Strict Mode

The AI must respond with `STATUS: PASSED` or `STATUS: FAILED` as the first line. If it doesn't:

1. Try Claude (most reliable at following instructions)
2. Check your rules file isn't confusing the AI
3. Temporarily disable strict mode: `STRICT_MODE="false"`

---

## Slow reviews on large files

The tool sends full file contents. For better performance:

```bash
# Add large/generated files to exclude
EXCLUDE_PATTERNS="*.min.js,*.bundle.js,dist/*,build/*,*.generated.ts"
```

---

## GitHub Models setup

```bash
# 1. Install GitHub CLI
brew install gh

# 2. Authenticate
gh auth login

# 3. Configure GGA
echo 'PROVIDER="github:gpt-4o"' > .gga

# Available models: https://github.com/marketplace/models
```

---

## Timeout issues

If reviews are timing out (exit code 124):

```bash
# Increase timeout (default: 300s)
TIMEOUT="600"          # In .gga config
GGA_TIMEOUT=600 gga run  # Or via environment variable

# Review fewer files at once
EXCLUDE_PATTERNS="*.min.js,*.bundle.js,dist/*"
```

---

## GGA not running from VS Code Source Control panel

If GGA doesn't trigger when committing from VS Code's Source Control UI:

1. Ensure the hook is installed: `ls -la .git/hooks/pre-commit`
2. Check that `gga` is in your PATH — VS Code may use a different shell profile
   - On Windows, check both PowerShell (`where gga`) and Git Bash (`which gga`) inside VS Code.
3. Try adding the full path in the hook:
   ```bash
   # .git/hooks/pre-commit
   /opt/homebrew/bin/gga run || exit 1
   ```
4. On Windows, if PATH still differs, hardcode the executable path in the hook (for example `C:/Users/<you>/.local/bin/gga.exe run || exit 1`).
5. Check the Git output channel (View → Output → Git) for error messages

---

## LM Studio connection issues

If you get "Failed to connect to LM Studio" errors:

1. Ensure LM Studio is running and the API server is enabled
2. Check the API port in LM Studio settings (default: 1234)
3. Verify the host setting:
   ```bash
   # Default
   LMSTUDIO_HOST="http://localhost:1234/v1"

   # Custom port
   LMSTUDIO_HOST="http://localhost:8080/v1"
   ```
4. Test the connection:
   ```bash
   curl http://localhost:1234/v1/models
   ```
