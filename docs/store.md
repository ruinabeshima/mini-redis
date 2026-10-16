# STORE Package

Defines a key-value store which is used in the Redis server.

## Architecture

The store is defined as a struct with a RWMutex lock and the key value store itself.

```go
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}
```

- The `New()` function acts as a constructor for the `Store` struct; it initialises the map in memory and returns a ready to use pointer to the struct
- The function `Set()` stores a key-value pair in the store. It grants access to only a single writer through `Lock()`, stores the key-value pair then unlocks through `Unlock`.
- The function `Get()` retrieves a value using a key. Since multiple readers can use this concurrently, `RLock()` and `RUnlock()` are used.
- `data` in this case is lowercase. Since it is unexported, no code outside of the package can reach the map at all, and only through the `Set()` and `Get()` methods (both of which lock using RWMutex).

## Usage

- Since the KVStore is defined as a struct with methods, it is constructed once in `main.go` before the accept loop, and shared by every connection goroutine for the lifetime of the server.
- The `HANDLER` package also utilises this store as one of its struct variables for the constructor.

## RWMutex

Originally, the store was defined as a simple variable storing a map of keys and values.

`var KVStore = make(map[string]string)`

However, this encountered a problem when multiple client goroutines tried to access and write into this KVStore at once. This data race can cause the Go runtime to crash as keys and values could be corrupted.

`sync.RWMutex` fixes this problem by allowing concurrent readers (`RLock`) but granting exclusive access to a single writer (`Lock`).

This gave rise to the idea of implementing the key-value store using Object Oriented Programming, with a constructor and methods `Set()` and `Get()`.
