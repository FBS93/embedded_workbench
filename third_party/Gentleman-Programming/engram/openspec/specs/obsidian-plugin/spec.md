# obsidian-plugin Specification

## Purpose

TypeScript Obsidian community plugin that provides a settings UI, a ribbon sync button, and HTTP-based incremental sync against the engram server, allowing users to trigger vault updates from inside Obsidian.

---

## Requirements

### Requirement: REQ-PLUGIN-01 — Plugin Manifest and Lifecycle

The plugin MUST include a valid `manifest.json` with `id`, `name`, `version`, `minAppVersion`, and `author` fields. The plugin MUST implement `onload()` and `onunload()` lifecycle hooks. On load, it MUST register its ribbon button, settings tab, and optional polling interval. On unload, it MUST clear any active polling interval.

#### Scenario: Plugin loads successfully

- GIVEN the user has Engram plugin installed in Obsidian
- WHEN Obsidian starts or the plugin is enabled
- THEN the plugin registers a ribbon button and a settings tab without errors
- AND no polling starts unless auto-sync interval > 0 is configured

#### Scenario: Plugin unloads cleanly

- GIVEN the plugin is active with a polling interval running
- WHEN the user disables the plugin or Obsidian closes
- THEN the polling interval is cleared and no further HTTP calls are made

---

### Requirement: REQ-PLUGIN-02 — Settings Tab

The plugin MUST provide a settings tab with the following fields:
- **Engram URL** (string, default: `http://localhost:4444`): base URL of the engram server
- **Sync interval** (number, default: 0): auto-sync interval in minutes; 0 means manual only
- **Vault subfolder** (string, default: `engram`): subfolder inside the vault where notes are written

Settings MUST be persisted via Obsidian's `loadData()` / `saveData()` API. Changes to sync interval MUST immediately restart or clear the polling interval.

#### Scenario: User changes sync interval from 0 to 5

- GIVEN sync interval is 0 (manual only)
- WHEN the user sets sync interval to 5 in the settings tab
- THEN a polling interval of 5 minutes is registered
- AND a sync triggers every 5 minutes until the interval is changed

#### Scenario: User clears Engram URL

- GIVEN the user deletes the Engram URL value and saves
- WHEN a sync is triggered
- THEN an error notice "Engram URL is required" is shown
- AND no HTTP request is made

---

### Requirement: REQ-PLUGIN-03 — Ribbon Button for Manual Sync

The plugin MUST add a ribbon button with a recognizable icon. Clicking it MUST trigger a manual sync against the configured Engram URL. While sync is in progress, the button MUST show a loading state. On completion, it MUST show a success or error notice.

#### Scenario: Successful manual sync

- GIVEN the plugin is configured with a valid Engram URL
- WHEN the user clicks the ribbon button
- THEN the button enters a loading state
- AND a sync request is sent to `{engram_url}/export`
- AND on success a notice "Synced N observations" is shown
- AND the button returns to idle state

#### Scenario: Server unreachable on manual sync

- GIVEN the Engram server is not running
- WHEN the user clicks the ribbon button
- THEN a notice "Sync failed: could not reach engram server" is shown
- AND the button returns to idle state

---

### Requirement: REQ-PLUGIN-04 — HTTP API Sync Mode

When a sync is triggered (manual or automatic), the plugin MUST:
1. Call `GET {engram_url}/export` (with optional query params: `project`, `since`)
2. Receive the export payload (list of notes with path, content, deleted flag)
3. Write new/updated notes to `{vault}/{subfolder}/{path}`
4. Delete vault files flagged as deleted in the payload
5. Update a local sync-state with the timestamp of the last successful sync

The plugin MUST NOT write files outside the configured vault subfolder.

#### Scenario: First sync — empty vault subfolder

- GIVEN the vault has no `engram/` subfolder
- WHEN a sync is triggered
- THEN the subfolder is created
- AND all exported notes are written to `{vault}/engram/`
- AND a success notice shows the count of written files

#### Scenario: Incremental sync — only new observations written

- GIVEN a previous sync was completed at T1
- WHEN a sync runs at T2 (T2 > T1)
- THEN the request includes `?since=T1` as a query parameter
- AND only observations created/updated after T1 are returned and written

#### Scenario: Deleted observation flagged in payload

- GIVEN the server marks obs-7 as deleted in the export payload
- WHEN the plugin processes the payload
- THEN `{vault}/engram/.../obs-7.md` is deleted from disk if it exists

---

### Requirement: REQ-PLUGIN-05 — Status Bar Indicator

The plugin MUST display a status bar item showing:
- Last sync timestamp (human-readable, e.g., "Last sync: 2 min ago")
- Count of observations on last sync (e.g., "42 notes")

The status bar item MUST update after every sync (success or failure). On failure, it MUST display "Sync failed" with the timestamp of the last attempt.

#### Scenario: Status bar after successful sync

- GIVEN a sync completes successfully with 37 observations
- WHEN the status bar updates
- THEN it displays "37 notes · synced just now" (or equivalent relative time)

#### Scenario: Status bar after failed sync

- GIVEN a sync attempt fails
- WHEN the status bar updates
- THEN it displays "Sync failed · {relative time of attempt}"
- AND the previously synced count is NOT overwritten
