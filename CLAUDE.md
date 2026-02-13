# sndit

Sound code generation tool and engine for Go games.

## Conventions

- Use `github.com/stretchr/testify/require` for all assertions
- Don't assert on specific error messages — only check `require.Error` / `require.NoError`
- TDD in small cycles: write one failing test, make it pass, then write the next test. Don't batch multiple tests before implementing.
- Extract testable functions from main (e.g. `run()`) so coverage is measured in-process
- When changing function signatures, update generated templates and their test expectations together
- Keep README.md up to date when changing user-facing behaviour (CLI flags, directory conventions, generated API)

## Commands

- `go test ./...`
- `go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out`
