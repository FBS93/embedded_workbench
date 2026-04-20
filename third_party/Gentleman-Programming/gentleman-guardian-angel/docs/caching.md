# Smart Caching

> 📖 Back to [README](../README.md)

GGA includes intelligent caching to speed up reviews by skipping files that haven't changed.

---

## How It Works

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cache Logic                               │
├─────────────────────────────────────────────────────────────────┤
│  1. Hash AGENTS.md + .gga config                                │
│     └─► If changed → Invalidate ALL cache                       │
│                                                                  │
│  2. For each staged file:                                        │
│     └─► Hash file content                                        │
│         └─► If hash exists in cache with PASSED → Skip          │
│         └─► If not cached → Send to AI for review               │
│                                                                  │
│  3. After PASSED review:                                         │
│     └─► Store file hash in cache                                │
└─────────────────────────────────────────────────────────────────┘
```

---

## Cache Invalidation

The cache automatically invalidates when:

| Change                | Effect                        |
| --------------------- | ----------------------------- |
| File content changes  | Only that file is re-reviewed |
| `AGENTS.md` changes   | **All files** are re-reviewed |
| `.gga` config changes | **All files** are re-reviewed |

---

## Cache Commands

```bash
# Check cache status
$ gga cache status

Cache Status:

  Cache directory: ~/.cache/gga/a1b2c3d4...
  Cache validity: Valid
  Cached files: 12
  Cache size: 4.0K

# Clear project cache
$ gga cache clear
✅ Cleared cache for current project

# Clear all cache (all projects)
$ gga cache clear-all
✅ Cleared all cache data
```

---

## Bypass Cache

```bash
# Force review all files, ignoring cache
gga run --no-cache
```

---

## Cache Location

```
~/.cache/gga/
├── <project-hash-1>/
│   ├── metadata          # Hash of AGENTS.md + .gga
│   └── files/
│       ├── <file-hash-a> # "PASSED"
│       └── <file-hash-b> # "PASSED"
└── <project-hash-2>/
    └── ...
```
