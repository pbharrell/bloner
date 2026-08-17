# Bloner Agent Guide

## Build And Run

- This is one Go 1.23.6 module (`github.com/pbharrell/bloner`) using Ebiten; `main.go` is the native and WASM entrypoint, and `graphics/` is its only internal package.
- Build the browser client with `./generate_wasm.sh`. It writes the checked-in `bloner.wasm`; rebuild it whenever client Go code changes before testing or serving the web build.
- Serve the browser build with `./deploy_webserver.sh` (`npx http-server ./ -o -p 8080`). `index.html` loads `wasm_exec.js` and `bloner.wasm` from the repository root.
- The game dials `ws://localhost:9000/ws` in `Game.init`; start a compatible `bloner-server` separately before exercising lobby or game flows.
- `./test_game.sh` is a manual multiplayer smoke script, not an automated test: it launches four native clients and leaves the first three running in the background. Do not use it as a clean CI-style check.

## Structure And Protocol

- Screen flow is a stack of `Page` implementations: the lobby page files create transitions, while `page_game_active.go` owns active-game UI and state.
- Keep client/server protocol changes coordinated with the external `github.com/pbharrell/bloner-server/connection` types. `message_handler.go` dispatches wire messages; `encode.go` and `decode.go` translate game state.
- Assets under `assets/` are embedded by `main.go` (`//go:embed assets/*`), so asset-path changes affect both native and WASM builds.

## Verification

- Run `go test ./...` for the available compile check. There are currently no Go test files.
- Run `gofmt -w *.go graphics/*.go` on edited Go files; no linter, formatter config, CI workflow, or task runner is present.
