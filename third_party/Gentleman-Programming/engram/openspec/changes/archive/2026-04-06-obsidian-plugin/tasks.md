# Tasks: Obsidian Plugin for Engram Brain Visualization

## Phase 1: Foundation — `internal/obsidian/` package

- [x] 1.1 **[RED]** Write `TestSlugify` in `internal/obsidian/slug_test.go` — covers empty title, unicode, long titles (>60 chars), collision by ID (REQ-EXPORT-07)
- [x] 1.2 **[GREEN]** Implement `Slugify(title string, id int64) string` in `internal/obsidian/slug.go`
- [x] 1.3 **[RED]** Write `TestObservationToMarkdown` in `internal/obsidian/markdown_test.go` — YAML frontmatter, body, wikilinks, no topic_key case (REQ-EXPORT-03)
- [x] 1.4 **[GREEN]** Implement `ObservationToMarkdown(obs store.Observation) string` in `internal/obsidian/markdown.go`
- [x] 1.5 **[RED]** Write `TestSyncStateRoundTrip` in `internal/obsidian/state_test.go` — marshal/unmarshal, missing file returns empty state (REQ-EXPORT-06)
- [x] 1.6 **[GREEN]** Implement `ReadState(path string) SyncState`, `WriteState(path string, s SyncState) error` in `internal/obsidian/state.go` with types `SyncState`, `ExportResult`
- [x] 1.7 **[RED]** Write `TestSessionHub` and `TestTopicHub` / `TestTopicHubSkipped` in `internal/obsidian/hub_test.go` — hub with backlinks, hub threshold ≥2, no hub for singleton (REQ-EXPORT-04, REQ-EXPORT-05)
- [x] 1.8 **[GREEN]** Implement `SessionHubMarkdown(sessionID string, obs []ObsRef) string` and `TopicHubMarkdown(prefix string, obs []ObsRef) string` in `internal/obsidian/hub.go`

## Phase 2: Exporter Engine — `internal/obsidian/exporter.go`

- [x] 2.1 **[RED]** Write `TestNewExporter` and `TestExportConfig` in `internal/obsidian/exporter_test.go` — missing vault path returns error, `StoreReader` mock satisfies interface
- [x] 2.2 **[GREEN]** Implement `ExportConfig`, `Exporter`, `StoreReader` interface, `NewExporter()`, and `Export()` skeleton in `internal/obsidian/exporter.go`
- [x] 2.3 **[RED]** Write `TestIncrementalExport` — state file gates which obs are written; only new/changed obs exported (REQ-EXPORT-06, REQ-EXPORT-08)
- [x] 2.4 **[GREEN]** Implement incremental filter logic (read state → filter by `LastExportAt` or `--since` → write only changed obs) in `Export()`
- [x] 2.5 **[RED]** Write `TestDeletedObsRemoved` — obs in state with `deleted_at != nil` triggers file deletion from vault (REQ-EXPORT-09)
- [x] 2.6 **[GREEN]** Implement deleted-obs cleanup pass in `Export()`
- [x] 2.7 **[RED]** Write `TestProjectFilter` — `--project` flag limits exported obs to matching project
- [x] 2.8 **[GREEN]** Implement project filter in `Export()` + `TestIdempotentExport` GREEN
- [x] 2.9 **[RED]** Write `TestFullExportPipeline` integration test (temp vault, in-memory store, 3 projects / 5 sessions / 20 obs / 1 deleted, assert dir structure + state file) (REQ-EXPORT-02)
- [x] 2.10 **[GREEN]** Fix any integration failures; ensure `Export()` creates all required subdirs

## Phase 3: CLI Command — `cmd/engram/main.go`

- [x] 3.1 **[RED]** Write CLI tests in `cmd/engram/main_test.go` — missing `--vault` exits with code 1, unknown flag exits 1, valid flags call injected exporter (REQ-EXPORT-01)
- [x] 3.2 **[GREEN]** Add `case "obsidian-export": cmdObsidianExport(cfg)` to switch; implement `cmdObsidianExport` with flag parsing and `var newObsidianExporter = obsidian.NewExporter` injection; update help text

## Phase 4: TypeScript Plugin — `plugin/obsidian/`

- [x] 4.1 Create `plugin/obsidian/manifest.json` (id, name, version, minAppVersion, author), `package.json`, `tsconfig.json`, `esbuild.config.mjs` — scaffolding only (REQ-PLUGIN-01)
- [x] 4.2 Implement `plugin/obsidian/src/settings.ts` — `EngramSettings` interface, defaults (`url`, `syncInterval`, `subfolder`), `EngramSettingTab` class with `loadData`/`saveData`; interval restart on change (REQ-PLUGIN-02)
- [x] 4.3 Implement `plugin/obsidian/src/main.ts` — `EngramBrainPlugin` class, `onload()` registers ribbon button + settings tab + status bar item + optional polling, `onunload()` clears interval (REQ-PLUGIN-01, REQ-PLUGIN-03, REQ-PLUGIN-05)
- [x] 4.4 Implement `plugin/obsidian/src/sync.ts` — `syncNow()`: `GET /export?since=T`, diff state, `vault.create/modify/delete`, write state file; guard against writes outside subfolder (REQ-PLUGIN-04)
- [x] 4.5 Wire status bar updates in `main.ts` — success: "N notes · synced just now", failure: "Sync failed · {relative time}" without overwriting previous count (REQ-PLUGIN-05)

## Phase 5: Polish

- [x] 5.1 Update `README.md` — add `obsidian-export` usage, flags table, vault structure diagram, and link to `plugin/obsidian/`
- [x] 5.2 Verify `goreleaser` config includes `plugin/obsidian/` assets; confirm `internal/obsidian/` compiles under `CGO_ENABLED=0`
