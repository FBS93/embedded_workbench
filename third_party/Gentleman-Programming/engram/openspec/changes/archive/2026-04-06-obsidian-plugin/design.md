# Technical Design: Obsidian Plugin for Engram Brain Visualization

## Executive Summary

Add Obsidian visualization to Engram via two deliverables: (1) a Go CLI exporter (`engram obsidian-export`) that reads SQLite and writes a markdown vault with frontmatter + wikilinks, and (2) an optional TypeScript Obsidian community plugin that provides a settings UI, ribbon button, and HTTP API sync mode.

---

## 1. Package Architecture

### New Packages

```
internal/obsidian/          ← Go export engine (Phase 1)
  exporter.go               ← Exporter struct, Export(), IncrementalExport()
  markdown.go               ← Observation → markdown conversion
  slug.go                   ← Title → filename-safe slug
  hub.go                    ← Session + topic hub generation
  state.go                  ← SyncState read/write
  exporter_test.go          ← Unit tests

plugin/obsidian/            ← TypeScript plugin (Phase 2)
  manifest.json
  src/
    main.ts                 ← Plugin class
    settings.ts             ← Settings tab
    sync.ts                 ← HTTP API sync client
  esbuild.config.mjs
  package.json
  tsconfig.json
```

### Dependency Graph (addition)

```
cmd/engram/main.go
  └── internal/obsidian  ← NEW (reads from store)
        └── internal/store  (read-only: Export(), AllObservations, RecentSessions, etc.)
```

`internal/obsidian` depends ONLY on `internal/store`. No reverse dependency. No new external Go deps.

### CLI Wiring (cmd/engram/main.go)

Follow existing pattern — add one case + one function:

```go
// In the switch statement (~line 138):
case "obsidian-export":
    cmdObsidianExport(cfg)

// Injectable for testing (with existing vars):
var newObsidianExporter = obsidian.NewExporter
```

Function signature follows `cmdExport`/`cmdSync` pattern: parse flags, open store, call exporter.

---

## 2. Core Types (Go)

```go
package obsidian

import "github.com/Gentleman-Programming/engram/internal/store"

// ExportConfig holds all CLI flags for the export command.
type ExportConfig struct {
    VaultPath string        // --vault (required)
    Project   string        // --project (optional filter)
    Limit     int           // --limit (0 = no limit)
    Since     time.Time     // --since (zero = use state file)
    Force     bool          // --force (ignore state, full re-export)
}

// Exporter reads from the store and writes markdown files to a vault.
type Exporter struct {
    store     StoreReader
    config    ExportConfig
}

// StoreReader is the read-only interface the exporter needs.
// Keeps the dependency narrow — easy to mock in tests.
type StoreReader interface {
    Export() (*store.ExportData, error)
    Stats() *store.Stats
}

// SyncState tracks what was exported and when.
type SyncState struct {
    LastExportAt string            `json:"last_export_at"`
    Files        map[int64]string  `json:"files"`         // obs ID → relative vault path
    SessionHubs  map[string]string `json:"session_hubs"`  // session ID → relative path
    TopicHubs    map[string]string `json:"topic_hubs"`    // topic prefix → relative path
    Version      int               `json:"version"`       // schema version (1)
}

// ExportResult summarizes what happened.
type ExportResult struct {
    Created    int
    Updated    int
    Deleted    int
    Skipped    int
    HubsCreated int
    Errors     []error
}

func NewExporter(s StoreReader, cfg ExportConfig) *Exporter
func (e *Exporter) Export() (*ExportResult, error)          // full or incremental
```

---

## 3. Data Transformation Pipeline

```
store.Export()
  → ExportData{Sessions[], Observations[], Prompts[]}
    → Filter by project (if --project set)
    → Filter by time (if --since or state file)
    → For each observation:
        1. Generate slug from title + id
        2. Determine folder: {vault}/engram/{project}/{type}/
        3. Build frontmatter (YAML)
        4. Build body (content as-is, it's already markdown)
        5. Append wikilinks: [[session-{session_id}]], [[topic-{prefix}]]
        6. Write file (create or overwrite)
    → For each session with ≥1 exported observation:
        1. Generate session hub note with backlinks
    → For each topic_key prefix with ≥2 observations:
        1. Generate topic hub note with backlinks
    → For deleted observations (in state but not in export):
        1. Remove vault file
        2. Remove from state
    → Write updated SyncState
```

---

## 4. Vault Structure

### Directory Layout

```
{vault}/
  engram/                              ← root namespace (MUST NOT write outside)
    .engram-sync-state.json            ← hidden state file
    engram/                            ← project folder
      architecture/
        exploration-obsidian-plugin-1719.md
        sdd-design-jwt-auth-1205.md
      decision/
        switched-sessions-to-jwt-892.md
      bugfix/
        fixed-fts5-syntax-error-1401.md
      _sessions/
        session-abc123.md              ← session hub
        session-def456.md
      _topics/
        sdd--obsidian-plugin.md        ← topic hub (for sdd/obsidian-plugin/*)
        architecture--auth-model.md    ← topic hub
    gentle-ai/                         ← another project
      ...
```

### Observation File Template

```markdown
---
id: 1719
type: architecture
project: engram
scope: project
topic_key: sdd/obsidian-plugin/explore
session_id: session-abc123
created_at: "2026-04-06T14:30:00Z"
updated_at: "2026-04-06T14:35:00Z"
revision_count: 1
tags:
  - engram
  - architecture
  - sdd
aliases:
  - "Exploration: Obsidian plugin for Engram brain visualization"
---

# Exploration: Obsidian plugin for Engram brain visualization

**What**: Complete exploration of Obsidian plugin architecture...
**Why**: User wants to visualize agent memory as connected graph...
**Where**: internal/server/, internal/store/...
**Learned**: Obsidian graph only shows wikilinks between real .md files...

---
*Session*: [[session-session-abc123]]
*Topic*: [[topic-sdd--obsidian-plugin]]
```

### Session Hub Template

```markdown
---
type: session-hub
session_id: session-abc123
project: engram
started_at: "2026-04-06T14:00:00Z"
ended_at: "2026-04-06T15:30:00Z"
tags:
  - session
  - engram
---

# Session: session-abc123

**Project**: engram | **Started**: 2026-04-06 14:00 | **Ended**: 15:30

## Summary
> Explored Obsidian plugin architecture and created SDD artifacts...

## Observations
- [[exploration-obsidian-plugin-1719]]
- [[sdd-proposal-obsidian-plugin-1720]]
- [[sdd-spec-obsidian-plugin-1721]]
```

### Topic Hub Template

```markdown
---
type: topic-hub
topic_prefix: sdd/obsidian-plugin
project: engram
tags:
  - topic
  - sdd
---

# Topic: sdd/obsidian-plugin

## Related Observations
- [[exploration-obsidian-plugin-1719]] (explore)
- [[sdd-proposal-obsidian-plugin-1720]] (proposal)
- [[sdd-spec-obsidian-plugin-1721]] (spec)
- [[sdd-design-obsidian-plugin-1722]] (design)
- [[sdd-tasks-obsidian-plugin-1723]] (tasks)
```

---

## 5. Slug Generation Algorithm

```go
func Slugify(title string, id int64) string {
    // 1. Lowercase
    s := strings.ToLower(title)
    // 2. Replace non-alphanumeric with hyphens
    s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
    // 3. Trim leading/trailing hyphens
    s = strings.Trim(s, "-")
    // 4. Truncate to 60 chars (avoid filesystem limits)
    if len(s) > 60 {
        s = s[:60]
        s = strings.TrimRight(s, "-")
    }
    // 5. Append ID for uniqueness
    return fmt.Sprintf("%s-%d", s, id)
}
```

Examples:
- `"Fixed FTS5 syntax error"` + id 1401 → `fixed-fts5-syntax-error-1401`
- `"SDD Proposal: obsidian-plugin"` + id 1720 → `sdd-proposal-obsidian-plugin-1720`
- `""` (empty title) + id 42 → `observation-42`

---

## 6. Incremental Sync Algorithm

```
1. Read SyncState from {vault}/engram/.engram-sync-state.json
   - If missing or --force: full export (state = empty)
   
2. Call store.Export() → full ExportData
   
3. Filter observations:
   - If --project: keep only matching project
   - If --since: keep only updated_at >= since
   - Else if state.LastExportAt != "": keep only updated_at >= LastExportAt
   - Else: keep all (first run)

4. For each kept observation:
   a. Generate expected path via Slugify + folder rules
   b. If path exists in state.Files AND content unchanged → skip
   c. If path exists but content changed → overwrite file, update state
   d. If path NOT in state → create file, add to state
   
5. For deleted observations:
   a. For each id in state.Files: if obs has DeletedAt != nil → remove file, remove from state
   
6. Generate hub notes:
   a. Session hubs: for each session with ≥1 exported obs
   b. Topic hubs: for each topic_key prefix with ≥2 obs (split on last "/")
   
7. Write updated SyncState with new LastExportAt = now()
```

**Conflict resolution**: Engram always wins. If user manually edits a vault file, next export overwrites it. This is a one-way mirror by design.

---

## 7. TypeScript Plugin Architecture (Phase 2)

### main.ts

```typescript
import { Plugin, PluginSettingTab, Notice } from 'obsidian';

interface EngramSettings {
    engramUrl: string;      // default: "http://127.0.0.1:7437"
    autoSyncMinutes: number; // default: 0 (disabled)
    projectFilter: string;   // default: "" (all)
}

export default class EngramBrainPlugin extends Plugin {
    settings: EngramSettings;
    syncInterval: number | null;
    
    async onload() {
        await this.loadSettings();
        this.addSettingTab(new EngramSettingTab(this.app, this));
        this.addRibbonIcon('brain', 'Sync Engram Brain', () => this.syncNow());
        this.addCommand({ id: 'sync-engram', name: 'Sync Engram Brain', callback: () => this.syncNow() });
        this.addStatusBarItem().setText('Engram: ready');
        if (this.settings.autoSyncMinutes > 0) this.startAutoSync();
    }
    
    async syncNow() {
        // 1. GET {engramUrl}/export
        // 2. Read .engram-sync-state.json from vault
        // 3. Diff: create/update/delete files
        // 4. Write updated state
        // 5. Show Notice with summary
    }
}
```

### HTTP Sync Flow

Same algorithm as Go CLI but in TypeScript:
- `fetch(settings.engramUrl + '/export')` → JSON
- Same diff logic against state file
- Uses Obsidian's `vault.create()` / `vault.modify()` / `vault.delete()` API
- Never touches files outside `engram/` folder

### Settings Tab

- Engram URL (text input, validated on save with `/health` check)
- Auto-sync interval (dropdown: disabled, 5min, 15min, 30min, 1hr)
- Project filter (text input, optional)
- Manual sync button with last sync timestamp

---

## 8. Testing Strategy

### Unit Tests (`internal/obsidian/`)

| Test | What It Verifies |
|------|-----------------|
| `TestSlugify` | Title → slug conversion, edge cases (empty, unicode, long) |
| `TestObservationToMarkdown` | Frontmatter + body + wikilinks generation |
| `TestSessionHub` | Hub note with correct backlinks |
| `TestTopicHub` | Hub generated only when ≥2 obs share prefix |
| `TestTopicHubSkipped` | No hub when only 1 obs has a prefix |
| `TestSyncStateRoundTrip` | JSON marshal/unmarshal of state file |
| `TestIncrementalExport` | Only new/changed obs are exported |
| `TestDeletedObsRemoved` | Soft-deleted obs → file removed from vault |
| `TestIdempotentExport` | Re-run produces identical files |
| `TestProjectFilter` | Only matching project observations exported |

### Integration Test

```go
func TestFullExportPipeline(t *testing.T) {
    // 1. Create temp dir as vault
    // 2. Create in-memory store with test fixtures
    // 3. Run Export()
    // 4. Walk vault dir → verify file count, structure, content
    // 5. Add more observations → run Export() again
    // 6. Verify only new files created, existing unchanged
}
```

### Test Fixtures

Use `store.New(cfg)` with temp dir (same pattern as `store_test.go`). Populate with:
- 3 projects, 5 sessions, 20 observations across types
- 2 topic_key clusters (≥2 obs each), 1 singleton (should not create hub)
- 1 soft-deleted observation

---

## 9. Key Decisions

| Decision | Choice | Alternative | Why |
|----------|--------|-------------|-----|
| Store access | `StoreReader` interface | Direct `*store.Store` | Testable without full DB; narrow dependency |
| Vault namespace | `{vault}/engram/` | `{vault}/` root | Security boundary; never clobber user notes |
| File naming | `{slug}-{id}.md` | UUID or hash | Human-readable + collision-safe |
| Hub threshold | ≥2 observations | Always create | Prevents orphan hub clutter |
| Topic prefix extraction | Split on last `/` | Full topic_key | Groups SDD phases: `sdd/change/*` |
| State file location | Inside vault | `~/.engram/` | Portable with vault; Git-trackable |
| Conflict resolution | Engram wins | Merge | One-way mirror = simple, predictable |
| Deleted obs | Remove file | Mark with frontmatter flag | Clean vault; deleted means deleted |
| Phase 2 transport | HTTP API (`/export`) | SQLite direct | Clean separation; no WAL lock conflicts |
| No new Go deps | String formatting for markdown | goldmark/template | Zero dependency increase; markdown is simple |
