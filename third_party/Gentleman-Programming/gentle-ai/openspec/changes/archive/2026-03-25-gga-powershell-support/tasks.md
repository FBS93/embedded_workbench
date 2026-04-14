# Tasks: GGA PowerShell Support

## Phase 1: Asset — Embed gga.ps1

- [x] 1.1 Create `internal/assets/gga/gga.ps1` with the static shim content from design (git-on-PATH lookup, `$args` forwarding, `exit $LASTEXITCODE`)
- [x] 1.2 Verify `internal/assets/gga/` embed directive picks up `.ps1` files (check `internal/assets/gga/assets.go` or equivalent `//go:embed` declaration; add `*.ps1` glob if missing)

## Phase 2: Core Implementation — runtime.go helpers

- [x] 2.1 Add `RuntimeBinDir(homeDir string) string` to `internal/components/gga/runtime.go` returning `~/.local/share/gga/bin`
- [x] 2.2 Add `RuntimePS1Path(homeDir string) string` to `internal/components/gga/runtime.go` returning `RuntimeBinDir(homeDir) + "/gga.ps1"`
- [x] 2.3 Add `EnsurePowerShellShim(homeDir string) error` to `internal/components/gga/runtime.go` using `assets.Read("gga/gga.ps1")` + `filemerge.WriteFileAtomic`

## Phase 3: Integration — call-site wiring

- [x] 3.1 In `internal/cli/sync.go:297`, after the `EnsureRuntimeAssets` call, add `if runtime.GOOS == "windows" { EnsurePowerShellShim(s.homeDir) }` guard
- [x] 3.2 In `internal/cli/run.go:506`, apply the same Windows-guarded `EnsurePowerShellShim(s.homeDir)` call after `EnsureRuntimeAssets`

## Phase 4: Testing

- [x] 4.1 In `internal/components/gga/runtime_test.go`, add `TestRuntimeBinDir` and `TestRuntimePS1Path` table tests asserting correct path suffixes
- [x] 4.2 Add `TestEnsurePowerShellShimCreatesFileWhenMissing`: `t.TempDir()` home, call once, assert file exists with embedded content
- [x] 4.3 Add `TestEnsurePowerShellShimOverwritesStaleShim`: write sentinel content, call, assert content replaced
- [x] 4.4 Add `TestEnsurePowerShellShimIsNoOpWhenContentMatches`: call twice, assert `ModTime` unchanged (mirrors `TestEnsureRuntimeAssetsIsNoOpWhenContentMatches`)
- [x] 4.5 Add `TestAssetGGAPS1IsEmbeddedAndReadable`: assert `assets.Read("gga/gga.ps1")` returns non-empty content

## Phase 5: Cleanup

- [x] 5.1 In `docs/platforms.md`, remove the Windows/PowerShell limitation note
