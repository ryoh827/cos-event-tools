# Repository Guidelines

Contributors build a tiny Go utility that reads cosplay event CSV exports and creates per-event directories. Use this guide for quick orientation and to keep changes consistent.

## Project Structure & Module Organization
- `go.mod` defines module `github.com/ryoh827/cos-event-tools` targeting Go 1.22.1.
- `cmd/cos-mkdir/main.go` contains the CLI entry point; keep additional commands under `cmd/<tool>/` with their own `main.go` files.
- Generated binaries such as `cos-mkdir` should be git-ignored or removed before committing; rebuild locally when needed.

## Build, Test, and Development Commands
- `go build ./cmd/cos-mkdir`: compile the CLI and surface syntax or type errors.
- `go run ./cmd/cos-mkdir data.csv`: execute the tool directly against a CSV export to verify end-to-end behavior.
- `go test ./...`: execute unit tests across future packages; currently returns immediately because no tests exist.

## Coding Style & Naming Conventions
- Format with `gofmt` (tabs for indentation, newline at EOF) and organize imports with `goimports` if available.
- Favor clear, exported names only when packages need to be reused; otherwise keep helpers unexported (lowerCamelCase) as in `parseDate`.
- Keep functions cohesive—`run` orchestrates I/O, helpers handle parsing or sanitizing; follow this pattern when adding features.

## Testing Guidelines
- Add table-driven tests in `_test.go` files alongside the package under test (e.g., `cmd/cos-mkdir/main_test.go`).
- Name tests `TestFunctionName_Scenario` for clarity, and mock filesystem interactions via temporary directories.
- Aim to cover edge cases such as malformed dates, missing columns, and existing directories; supplement with manual `go run` checks using sample CSVs.

## Commit & Pull Request Guidelines
- The repo history is empty, so establish clear, imperative commit messages (e.g., `feat: add honorific toggle`).
- Reference issues in the body when applicable and describe observable behavior changes plus validation steps.
- Pull requests should summarize intent, list testing evidence (`go build`, `go test`, manual CSV run), and include screenshots or sample output when UX-visible.
