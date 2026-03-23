# AGENTS.md - Coding Agent Instructions

## Project Overview

Go module (`github.com/appscode/static-assets`) providing static assets (JSON data, images, files) for AppsCode products. Uses Go embed for serving data as `fs.FS`. Module path: `appscode-cloud`. Go version: 1.21+ (CI uses 1.25).

## Build & Test Commands

```bash
# Build all packages
CGO_ENABLED=0 go build -v ./...

# Run all tests
CGO_ENABLED=0 go test -v ./...

# Run a single test
CGO_ENABLED=0 go test -v -run TestName ./path/to/package/...

# Run go vet
go vet ./...

# Format code (uses vendored mod)
hack/fmt.sh api cmd data *.go
# Or manually:
goimports -w <files>
gofmt -s -w <files>

# Publish static assets to S3
make publish
```

All build/test commands require `CGO_ENABLED=0` and `GOFLAGS="-mod=vendor"`.

## Code Style

### Package Names
- Lowercase, single-word: `api`, `staticassets`
- Root package name is `staticassets` (not `static-assets`)

### Imports
- Group stdlib first, then external packages, separated by blank line
- Use vendored dependencies (`-mod=vendor`)
- Import the root module as: `staticassets "github.com/appscode/static-assets"`
- Subpackages: `"github.com/appscode/static-assets/api"`

```go
import (
	"encoding/json"
	"os"

	"github.com/appscode/static-assets/api"
)
```

### Naming
- Exported types/funcs: PascalCase (`Product`, `DownloadReleaseAssets`)
- Unexported: camelCase (`exists`, `repoList`)
- Acronyms kept uppercase: `API`, `URL`, `FS`, `ID` (e.g., `URLRef`, `RepoURL`)
- File names: `snake_case.go` for product data files, `main.go` for commands

### Struct Tags
- Always include JSON tags on struct fields
- Use `omitempty` for optional fields
- Prefer explicit JSON field names: `json:"fieldName"`

### Error Handling
- Use `panic()` only for unrecoverable programmer errors (e.g., failed embed.Sub)
- Return `error` for runtime/IO failures; handle or propagate
- Use `log.Fatalln()` for fatal errors in CLI tools
- Use type-switch on error types for specific handling (see GitHub API errors)

### Embedding
- Use `//go:embed` directives for static data
- Wrap in `fs.FS` via `fs.Sub()` for clean access paths

### General
- Prefer `bytes.Buffer` + `json.NewEncoder` over `json.Marshal` for formatted JSON output
- Use `os.WriteFile` with `0o644` permissions
- Use `0o755` for directories
- Keep functions small and focused
- No comments unless explicitly requested
