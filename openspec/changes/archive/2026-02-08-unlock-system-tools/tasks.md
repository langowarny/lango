## 1. Filesystem Tools Integration

- [x] 1.1 Register `fs_list` tool handler in `internal/app/app.go` using `fsTool.ListDir`.
- [x] 1.2 Register `fs_write` tool handler in `internal/app/app.go` using `fsTool.Write`.
- [x] 1.3 Register `fs_edit` tool handler in `internal/app/app.go` using `fsTool.Edit`.
- [x] 1.4 Register `fs_mkdir` tool handler in `internal/app/app.go` using `fsTool.Mkdir`.
- [x] 1.5 Register `fs_delete` tool handler in `internal/app/app.go` using `fsTool.Delete`.

## 2. Execution Tools Integration

- [x] 2.1 Refactor existing `exec` tool handler in `internal/app/app.go` to remove redundant error check.
- [x] 2.2 Register `exec_bg` tool handler in `internal/app/app.go` using `sv.execTool.StartBackground`.
- [x] 2.3 Register `exec_status` tool handler in `internal/app/app.go` using `sv.execTool.GetBackgroundStatus`.
- [x] 2.4 Register `exec_stop` tool handler in `internal/app/app.go` using `sv.execTool.StopBackground`.

## 3. Verification

- [x] 3.1 Create a verification script to test new filesystem capabilities (listing, creating, editing, deleting).
- [x] 3.2 Create a verification script to test background execution (starting a simple loop, checking status, stopping it).
- [x] 3.3 Verify all new tools are listed and functional in a live `lango serve` session.
