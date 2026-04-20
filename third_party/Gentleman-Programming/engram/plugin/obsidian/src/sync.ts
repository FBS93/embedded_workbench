import { Notice, Vault } from "obsidian";
import type EngramBrainPlugin from "./main";

// ─── API Types ─────────────────────────────────────────────────────────────────

/** A single note entry returned by GET /export */
interface ExportNote {
	/** Relative path inside the vault subfolder, e.g. "engram/architecture/my-note-42.md" */
	path: string;
	/** Full markdown content of the note */
	content: string;
	/** If true, this note should be deleted from the vault */
	deleted: boolean;
}

/** Full response body from GET /export */
interface ExportData {
	notes: ExportNote[];
	/** ISO timestamp representing the export point (use as `since` on next sync) */
	exported_at: string;
	/** Total number of notes in this export */
	count: number;
}

// ─── Sync State ───────────────────────────────────────────────────────────────

interface SyncState {
	/** ISO timestamp of the last successful sync — used as `?since=` on the next request */
	lastSyncAt: string;
	/** Map of relative path → content hash (sha256 of content) — for change detection */
	files: Record<string, string>;
	/** Schema version */
	version: number;
}

const STATE_VERSION = 1;
const STATE_FILENAME = ".engram-sync-state.json";

// ─── Result ───────────────────────────────────────────────────────────────────

export interface SyncResult {
	created: number;
	updated: number;
	deleted: number;
	skipped: number;
	total: number;
}

// ─── Main Sync Function ───────────────────────────────────────────────────────

/**
 * Execute a full or incremental sync against the Engram server.
 *
 * Steps:
 * 1. Validate settings (URL required)
 * 2. GET {engramUrl}/export?since=T&project=P
 * 3. Read .engram-sync-state.json from vault subfolder
 * 4. Diff: create new, update changed, delete removed
 * 5. Write updated state file
 * 6. Return SyncResult for status bar + Notice
 *
 * SAFETY: Never touches files outside the configured vault subfolder.
 */
export async function syncNow(plugin: EngramBrainPlugin): Promise<SyncResult> {
	const { settings } = plugin;
	const vault: Vault = plugin.app.vault;
	const subfolder = settings.vaultSubfolder || "engram";

	// ── 1. Validate ──────────────────────────────────────────────────────────
	if (!settings.engramUrl) {
		new Notice("Engram URL is required");
		throw new Error("Engram URL is required");
	}

	// ── 2. Fetch export from server ──────────────────────────────────────────
	const exportData = await fetchExport(settings.engramUrl, {
		project: settings.projectFilter || undefined,
		since: settings.lastSyncAt || undefined,
	});

	// ── 3. Read current sync state ────────────────────────────────────────────
	const stateFile = `${subfolder}/${STATE_FILENAME}`;
	const state = await readState(vault, stateFile);

	// ── 4. Diff and write files ───────────────────────────────────────────────
	const result: SyncResult = {
		created: 0,
		updated: 0,
		deleted: 0,
		skipped: 0,
		total: exportData.count,
	};

	// Ensure the subfolder exists
	await ensureFolder(vault, subfolder);

	for (const note of exportData.notes) {
		// SAFETY: guard — path must stay inside subfolder
		const safePath = toSafePath(subfolder, note.path);
		if (!safePath) {
			console.warn(`[Engram] Skipping unsafe path: ${note.path}`);
			result.skipped++;
			continue;
		}

		if (note.deleted) {
			// Delete the file if it exists in the vault
			const file = vault.getFileByPath(safePath);
			if (file) {
				await vault.delete(file);
				delete state.files[safePath];
				result.deleted++;
			} else {
				result.skipped++;
			}
			continue;
		}

		// Check if content is unchanged (avoid unnecessary writes)
		const existingHash = state.files[safePath];
		const newHash = await hashContent(note.content);

		if (existingHash === newHash) {
			result.skipped++;
			continue;
		}

		// Ensure parent folder exists
		const parentFolder = safePath.substring(0, safePath.lastIndexOf("/"));
		if (parentFolder && parentFolder !== subfolder) {
			await ensureFolder(vault, parentFolder);
		}

		const existingFile = vault.getFileByPath(safePath);
		if (existingFile) {
			await vault.modify(existingFile, note.content);
			result.updated++;
		} else {
			await vault.create(safePath, note.content);
			result.created++;
		}

		state.files[safePath] = newHash;
	}

	// ── 5. Write updated state ────────────────────────────────────────────────
	state.lastSyncAt = exportData.exported_at || new Date().toISOString();
	state.version = STATE_VERSION;
	await writeState(vault, stateFile, state);

	return result;
}

// ─── HTTP ─────────────────────────────────────────────────────────────────────

async function fetchExport(
	baseUrl: string,
	params: { project?: string; since?: string }
): Promise<ExportData> {
	const url = new URL(`${baseUrl}/export`);
	if (params.project) url.searchParams.set("project", params.project);
	if (params.since) url.searchParams.set("since", params.since);

	let res: Response;
	try {
		res = await fetch(url.toString(), {
			signal: AbortSignal.timeout(30_000),
		});
	} catch (err) {
		throw new Error("Sync failed: could not reach engram server");
	}

	if (!res.ok) {
		throw new Error(
			`Sync failed: server returned ${res.status} ${res.statusText}`
		);
	}

	const data = (await res.json()) as ExportData;

	// Normalise: ensure `notes` is always an array
	if (!Array.isArray(data.notes)) {
		data.notes = [];
	}

	return data;
}

// ─── State File ───────────────────────────────────────────────────────────────

async function readState(vault: Vault, statePath: string): Promise<SyncState> {
	const empty: SyncState = {
		lastSyncAt: "",
		files: {},
		version: STATE_VERSION,
	};

	try {
		const exists = await vault.adapter.exists(statePath);
		if (!exists) return empty;

		const raw = await vault.adapter.read(statePath);
		const parsed = JSON.parse(raw) as Partial<SyncState>;
		return {
			lastSyncAt: parsed.lastSyncAt ?? "",
			files: parsed.files ?? {},
			version: parsed.version ?? STATE_VERSION,
		};
	} catch {
		return empty;
	}
}

async function writeState(
	vault: Vault,
	statePath: string,
	state: SyncState
): Promise<void> {
	const content = JSON.stringify(state, null, 2);
	const exists = await vault.adapter.exists(statePath);
	if (exists) {
		const file = vault.getFileByPath(statePath);
		if (file) {
			await vault.modify(file, content);
			return;
		}
	}
	await vault.create(statePath, content);
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

/**
 * Ensure a folder exists in the vault, creating it recursively if needed.
 */
async function ensureFolder(vault: Vault, folderPath: string): Promise<void> {
	const exists = await vault.adapter.exists(folderPath);
	if (!exists) {
		await vault.createFolder(folderPath);
	}
}

/**
 * Convert a server-returned note path to a safe vault path.
 *
 * Returns null if the resulting path would escape the subfolder.
 * This is the SAFETY GUARD — nothing may be written outside `subfolder/`.
 */
function toSafePath(subfolder: string, notePath: string): string | null {
	// Strip any leading slashes or dots from the server-provided path
	const cleaned = notePath.replace(/^[./\\]+/, "").replace(/\\/g, "/");

	// If the path already starts with the subfolder, keep it as-is
	let fullPath: string;
	if (cleaned.startsWith(`${subfolder}/`)) {
		fullPath = cleaned;
	} else {
		fullPath = `${subfolder}/${cleaned}`;
	}

	// Guard: resolved path must start with subfolder/
	if (!fullPath.startsWith(`${subfolder}/`)) {
		return null;
	}

	// Guard: reject path traversal attempts
	if (fullPath.includes("..")) {
		return null;
	}

	return fullPath;
}

/**
 * Produce a deterministic hash of note content for change detection.
 * Uses Web Crypto API (available in Obsidian's Electron environment).
 */
async function hashContent(content: string): Promise<string> {
	const encoded = new TextEncoder().encode(content);
	const hashBuffer = await crypto.subtle.digest("SHA-256", encoded);
	const hashArray = Array.from(new Uint8Array(hashBuffer));
	return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}
