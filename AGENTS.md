# memkv — agent instructions

Single Go module (`github.com/beeleelee/memkv`), no external dependencies.

## Commands

```sh
go test ./...          # all tests
go test -v -run TestPutGetDelete  # single test
go test -bench=. ./... # benchmarks
```

## Architecture

- `memkv.go` — package root, defines `ErrNotFound`
- `interface.go` — `KVDB` interface (Put/Delete/Contains/Get/Find/Num) + helpers (`BytesCompare`, `BytesClone`, `NoCopy`)
- `skiplist.go` — skip-list implementation with per-level heads, `sync.RWMutex` for concurrency, seeded deterministic rand (`0xdeadbeaf`)
- `NewBytesKV()` is the public constructor — returns a `KVDB` with byte-slice keys and value cloning

## Conventions

- Commits use `type/scope` prefix (e.g. `refine/`, `test/`, `add/`, `chore/`)
- No external test framework — uses only `testing.T` and `log.Fatal`
- Value copying is pluggable via `CopyValueFunc`; `BytesClone` is the default for byte keys
- `skiplist.go` has an unexported `print()` method for debugging link structure
