# memkv

An in-memory key-value store for Go, backed by a skip list with concurrent read/write support.

## Features

- **KVDB interface** — Put, Get, Delete, Contains, Find, Num
- **Concurrent safe** — `sync.RWMutex`; reads are shared, writes are exclusive
- **Pluggable comparison** — any key type via `CompareFunc`
- **Pluggable value cloning** — control whether values are copied on insert/read
- **Byte-slice convenience** — `NewBytesKV()` returns a ready-to-use `KVDB` with `[]byte` keys
- **No dependencies** — standard library only

## Install

```sh
go get github.com/beeleelee/memkv
```

## Quick start

```go
package main

import (
    "fmt"
    "github.com/beeleelee/memkv"
)

func main() {
    db := memkv.NewBytesKV()

    db.Put([]byte("name"), []byte("memkv"))
    db.Put([]byte("lang"), []byte("go"))

    val, err := db.Get([]byte("name"))
    if err != nil {
        panic(err)
    }
    fmt.Println(string(val.([]byte))) // memkv

    fmt.Println(db.Contains([]byte("lang"))) // true

    db.Delete([]byte("lang"))

    fmt.Println(db.Num()) // 1
}
```

## API

### `KVDB` interface

```go
type KVDB interface {
    Put(key, value any) error
    Delete(key any) error
    Contains(key any) bool
    Get(key any) (value any, err error)
    Find(key any) (rkey, value any, err error)
    Num() int
}
```

| Method | Description |
|--------|-------------|
| `Put` | Insert or update a key-value pair. On update the new value is copied. |
| `Get` | Retrieve a value by exact key. Returns `ErrNotFound` if missing. |
| `Delete` | Remove a key. Returns `ErrNotFound` if missing. |
| `Contains` | Check whether a key exists (no value copy). |
| `Find` | Return the first key-value pair ≥ the given key (useful for range queries / prefix scans). Returns `ErrNotFound` when no greater key exists. |
| `Num` | Return the current number of entries. |

### Constructor

```go
func NewBytesKV() KVDB
```

Creates a skip-list with `BytesCompare` (lexicographic `[]byte` order) and `BytesClone` (copy on read/write).

### Custom key types

```go
import "github.com/beeleelee/memkv"

sl := memkv.NewSkipList(myCompare, myCopy)
```

Where `myCompare` satisfies `CompareFunc` and `myCopy` satisfies `CopyValueFunc`.

### Helpers

| Function | Signature | Purpose |
|----------|-----------|---------|
| `BytesCompare` | `func(a, b any) int` | Compares `a.([]byte)` and `b.([]byte)` |
| `BytesClone` | `func(v any) any` | Returns a byte-for-byte copy of `v.([]byte)` |
| `NoCopy` | `func(v any) any` | Returns the same value without copying |

### Sentinel errors

```go
var ErrNotFound = errors.New("Key not found")
```

## Tests

```sh
go test ./...           # unit tests
go test -v -run TestPut # single test
go test -bench=. ./...  # benchmarks
```

## License

MIT
