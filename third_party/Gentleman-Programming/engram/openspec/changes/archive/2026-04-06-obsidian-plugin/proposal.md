# Proposal: Obsidian Brain Visualization

## Intent

Users have no way to visualize how their AI agent's memories connect. Engram stores rich relational data (topic_key hierarchies, session grouping, project isolation, type clustering) but it's invisible. Obsidian's native graph view can surface these connections as an interactive knowledge graph — if we write real `.md` files with wikilinks.

## Scope

### In Scope
- **`engram obsidian-export` CLI command** — new subcommand on existing binary
- **Markdown vault generation** — observations as notes with YAML frontmatter + wikilinks
- **Folder structure**: `{vault}/engram/{project}/{type}/{observation-slug}.md`
- **Hub notes** — per-session and per-topic-cluster aggregation pages
- **Incremental sync** — track last export timestamp, only write new/changed observations
- **TypeScript Obsidian plugin** (Phase 2) — settings UI, ribbon sync button, HTTP API polling

### Out of Scope
- Embeddings / RAG / semantic search
- Custom graph renderer (use Obsidian native)
- Bidirectional sync (Obsidian → Engram writes)
- Mobile Obsidian support
- Dataview queries (users can add their own)

## Capabilities

### New Capabilities
- `obsidian-export`: CLI command that reads SQLite store and writes a structured markdown vault with frontmatter, wikilinks, and hub notes for Obsidian graph visualization
- `obsidian-plugin`: TypeScript Obsidian community plugin providing settings UI, ribbon sync, and HTTP-based incremental sync against the engram server

### Modified Capabilities
- None — this is purely additive. Existing CLI, store, and server remain untouched.

## Approach

**Phase 1 — CLI exporter (MVP, ~3-5 days)**:
- New `cmdObsidianExport(cfg)` in `cmd/engram/main.go` switch statement
- Export logic in `internal/obsidian/` package: `Exporter` struct takes `*store.Store`, vault path, options
- Reads via `store.Export()` (full dump) or `store.RecentObservations()` (incremental)
- Writes markdown files: frontmatter (type, project, scope, topic_key, session_id, created_at) + content body + wikilinks section
- Wikilinks derived from: shared topic_key prefix → `[[topic-hub]]`, session_id → `[[session-hub]]`, project → folder
- Hub notes: `_sessions/{session-id}.md` lists all observations in that session; `_topics/{topic-prefix}.md` lists all observations sharing a topic_key prefix
- Incremental: store last-export timestamp in `{vault}/engram/.engram-sync-state.json`
- No new Go dependencies — markdown is string formatting, YAML frontmatter is `fmt.Sprintf`

**Phase 2 — Obsidian plugin (~5-7 days)**:
- Source in `plugin/obsidian/` following existing plugin directory pattern
- Standard Obsidian plugin boilerplate: `manifest.json`, `main.ts`, `esbuild.config.mjs`
- Settings tab: engram server URL, auto-sync interval, vault subfolder
- Ribbon button triggers `GET /export` → diffs against existing vault → writes/updates files
- Optional polling mode for live sync during coding sessions

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `cmd/engram/main.go` | Modified | Add `obsidian-export` case + `cmdObsidianExport()` function |
| `internal/obsidian/` | New | Export logic: `Exporter`, markdown rendering, wikilink generation, incremental state |
| `internal/obsidian/exporter_test.go` | New | Unit tests for markdown generation, wikilink extraction, incremental logic |
| `plugin/obsidian/` | New | TypeScript Obsidian community plugin (Phase 2) |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Large vaults (10k+ observations) slow Obsidian | Medium | Folder-per-project isolation, incremental sync, configurable project filter |
| Topic-key wikilinks create orphan hub pages | Low | Only generate hub notes when ≥2 observations share a prefix |
| Filename collisions from title slugification | Low | Append observation ID to slug: `{slug}-{id}.md` |
| Obsidian plugin API breaking changes | Low | Pin API version in manifest.json, Phase 2 is optional |

## Rollback Plan

- CLI command is a new case in the switch — remove the `case "obsidian-export"` line and delete `internal/obsidian/`
- Plugin is entirely in `plugin/obsidian/` — delete the directory
- No existing code is modified beyond adding one switch case
- Generated vault files are in a subfolder — users delete `{vault}/engram/`

## Dependencies

- None for Phase 1 (pure Go, no new imports)
- Phase 2: `obsidian` npm types package, `esbuild` for bundling (devDependencies only)

## Success Criteria

- [ ] `engram obsidian-export --vault ~/my-vault` generates a browsable Obsidian vault
- [ ] Obsidian graph view shows connected nodes via wikilinks (sessions, topics, projects)
- [ ] Incremental re-export only writes new/changed files (idempotent)
- [ ] Unit tests cover markdown generation, slug creation, wikilink extraction, and incremental state
- [ ] Phase 2: Obsidian plugin syncs via HTTP API with settings UI and ribbon button
