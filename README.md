# emptycasevet

A tiny `go/analysis` vet tool that reports empty non-default `case` clauses in `switch` / `type switch`.

If you intentionally leave a case empty, add at least a comment inside the case body to pass the check.

## Install
```bash
# Requires Go 1.24+
go install github.com/ivanpointer/emptycasevet/cmd/emptycasevet@latest
```

## Usage
- As a standalone analyzer:
```bash
emptycasevet ./...
```

- With `go vet` via `-vettool`:
```bash
emptycasevet_path=$(command -v emptycasevet)
GOFLAGS="-vettool=$emptycasevet_path" go vet ./...
```

- In CI pipelines, add a step (example GitHub Actions):
```yaml
- name: Install emptycasevet
  run: go install github.com/ivanpointer/emptycasevet/cmd/emptycasevet@latest
- name: Run emptycasevet
  run: emptycasevet ./...
```

## What it flags
- Non-default `case` with empty body: flagged.
- `default` with empty body: allowed.
- Non-default `case` with only a comment: allowed.

Example:
```go
switch x := 1; x {
case 1: // flagged: empty case body; did you mean `case a, b:`?
case 2:
    // ok with comment
default:
    // empty default allowed
}
```


## Library usage (for integrators)
If you want to consume the analyzer programmatically (e.g. for integration in a tool such as golangci-lint), import the package and reference the exported `Analyzer`:

```go
import emptycasevet "github.com/ivanpointer/emptycasevet"

func init() {
    // register emptycasevet.Analyzer with your driver/runner
    _ = emptycasevet.Analyzer
}
```

## Configuration options
The analyzer supports the following option (via go/analysis flags):

- `-allow_header_comment` (default: false):
  Consider an inline comment placed on the same line as the `case` header as acceptable.
  By default, only comments in the case body or on the next line(s) after the colon are considered.

Examples:
- Standalone:
  ```bash
  emptycasevet -allow_header_comment ./...
  ```
- With `go vet`:
  ```bash
  emptycasevet_path=$(command -v emptycasevet)
  GOFLAGS="-vettool=$emptycasevet_path -emptycase.allow_header_comment" go vet ./...
  ```
  Note: When used as a `-vettool`, flags are namespaced by the analyzer name (`emptycase`).

## Notes for adding to golangci-lint
- This repository exposes a public `Analyzer` at package path `github.com/ivanpointer/emptycasevet`.
- The CLI binary remains available at `github.com/ivanpointer/emptycasevet/cmd/emptycasevet`.
- The analyzer is small, has tests via `analysistest`, and is licensed under the MIT License (see LICENSE).
- When submitting a PR to golangci-lint, wire it similarly to other go/analysis analyzers by importing the package and registering the Analyzer.

## Stability, performance, and maintenance
- Stable: small, focused analyzer with a minimal API surface (single exported `Analyzer`).
- Performance: single-pass AST walk over switch/type switch nodes; no type-checker loading nor heavy allocations.
- Maintenance: unit tests cover both positive and negative cases; semantic versioning will be used for releases.


## Why this exists
Accidentally leaving a non-default case empty is easy when refactoring a switch (e.g., splitting values across cases or removing code). This linter nudges you to either:
- Merge values into a single case (e.g., case a, b:), or
- Leave an intentional comment in the case body to document that it's deliberately empty.

## Type switch example
```go
type T interface{}

func g(v T) {
    switch v.(type) {
    case int: // flagged: empty case body; did you mean `case a, b:`?
    case string:
        // ok: intentional, documented
    default:
        // default may be empty
    }
}
```

## Exit status
- Exits with code 0 if no issues are found
- Exits with code >0 (non-zero) if any diagnostics are reported
This makes it safe to use in CI to fail builds when empty cases are detected.

## Go compatibility
- Requires Go 1.24+
- Developed and tested with Go 1.24.x
- Uses the standard go/analysis API and does not rely on unstable internals

## Contributing
Contributions are welcome! Please:
- Open an issue describing the improvement or bug.
- Include tests for behavior changes.
- Keep the analyzer minimal and focused to reduce false positives and runtime overhead.

## Releases
Stable releases are tagged using semantic versioning (v0.x.y initially). Use a tagged version in CI for reproducible builds:
```bash
go install github.com/ivanpointer/emptycasevet/cmd/emptycasevet@v0.1.0
```
If you don’t need pinning, @latest is fine for local use.
