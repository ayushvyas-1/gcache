# GCache 🚀

A high-performance, thread-safe LRU (Least Recently Used) cache implementation in Go with both in-memory and network server capabilities.

## ✨ Features

- **Thread-Safe**: Full concurrent read/write support with mutex locking
- **LRU Eviction**: Automatic eviction of least recently used items when capacity is reached
- **TCP Server**: Network-accessible cache server with Redis-like protocol
- **Interactive Client**: Command-line client with interactive mode
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **High Performance**: Optimized with doubly linked list and hash map combination
- **Comprehensive Testing**: Unit tests, benchmarks, and concurrency tests included

## 🏗️ Architecture

GCache uses a classic LRU implementation combining:
- **HashMap**: O(1) key lookup
- **Doubly Linked List**: O(1) insertion/deletion and LRU ordering
- **Mutex Locks**: Thread-safe concurrent access

```
┌─────────────┐    ┌──────────────────┐
│   HashMap   │───▶│ DoublyLinkedList │
│  key->node  │    │  MRU ←→...←→ LRU │
└─────────────┘    └──────────────────┘
```

## 🚀 Quick Start

### Build the Project
```bash
make build
```

### Start the Server
```bash
# Start server with default settings (localhost:8080, capacity: 1000)
make run-server

# Or with custom settings
./gcache -mode=server -addr=localhost:9000 -capacity=5000
```

### Use the Interactive Client
```bash
make run-client
```

### Run a Quick Demo
```bash
make demo
```

## 📖 Usage

### In-Memory Cache

```go
package main

import (
    "fmt"
    "github.com/ayushvyas-1/gcache/internal/cache"
)

func main() {
    // Create cache with capacity of 100
    lru := cache.NewLRUCache(100)
    
    // Set values
    lru.Put("user:1", "John Doe")
    lru.Put("user:2", "Jane Smith")
    
    // Get values
    if value, exists := lru.Get("user:1"); exists {
        fmt.Printf("Found: %s\n", value)
    }
    
    // Check cache state
    fmt.Printf("Cache size: %d\n", lru.Size())
    lru.Print() // Debug output
}
```

### TCP Server Protocol

The server implements a Redis-like text protocol:

#### Commands

| Command | Syntax | Description | Response |
|---------|--------|-------------|----------|
| **GET** | `GET key` | Retrieve value for key | `+value` or `-ERR key not found` |
| **SET** | `SET key value` | Store key-value pair | `+OK` |
| **DEL** | `DEL key` | Delete key | `+OK` or `-ERR key not found` |
| **SIZE** | `SIZE` | Get cache size | `:number` |
| **CLEAR** | `CLEAR` | Clear all items | `+OK` |
| **PING** | `PING [message]` | Ping server | `+PONG` or `+message` |
| **INFO** | `INFO` | Server information | `+info_string` |
| **STATS** | `STATS` | Cache statistics | `+stats_string` |

#### Response Format
- `+OK` - Success response
- `+value` - String value response  
- `:number` - Integer response
- `-ERR message` - Error response

### Client Examples

#### Command Line Usage
```bash
# Single command
./gcache -mode=client -cmd "SET mykey hello"
./gcache -mode=client -cmd "GET mykey"

# Interactive mode
./gcache -mode=client -interactive
gcache> SET user:1 "John Doe"
OK: 
gcache> GET user:1
OK: John Doe
gcache> SIZE
VALUE: 1
```

#### Programmatic Client
```go
client, err := cache.NewClient("localhost:8080")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Set a value
err = client.Set("mykey", "myvalue")

// Get a value
value, err := client.Get("mykey")

// Delete a key
err = client.Delete("mykey")

// Get cache size
size, err := client.Size()
```

## 🔧 Configuration

### Server Options
```bash
./gcache -mode=server \
         -addr=localhost:8080 \    # Server address
         -capacity=1000            # Cache capacity
```

### Client Options
```bash
./gcache -mode=client \
         -addr=localhost:8080 \    # Server address
         -interactive \            # Interactive mode
         -cmd="GET mykey"          # Single command
```

## 🧪 Testing

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Benchmarks
```bash
go test ./tests/... -bench=. -benchmem
```

### Run Race Detection
```bash
go test ./tests/... -race
```

### Example Benchmark Results
```
BenchmarkPut-8           5000000    250 ns/op    48 B/op    2 allocs/op
BenchmarkGet-8          10000000    150 ns/op     0 B/op    0 allocs/op
BenchmarkConcurrent-8    2000000    800 ns/op    48 B/op    2 allocs/op
```

## 📁 Project Structure

```
gcache/
├── cmd/server/           # Main application entry point
│   └── main.go
├── internal/cache/       # Core cache implementation
│   ├── cache.go         # LRU cache logic
│   ├── list.go          # Doubly linked list
│   ├── node.go          # List node implementation
│   ├── TCP_Server.go    # TCP server
│   └── TCP_Client.go    # TCP client
├── tests/               # Test suite
│   └── cache_test.go
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── README.md           # This file
```
## ⚡ Performance Characteristics

- **Time Complexity**:
  - GET: O(1)
  - PUT: O(1) 
  - DELETE: O(1)
- **Space Complexity**: O(capacity)
- **Concurrency**: Full thread-safety with RWMutex
- **Memory**: Minimal overhead with efficient data structures

## 🛡️ Thread Safety

GCache is fully thread-safe:
- Uses `sync.RWMutex` for concurrent access
- Read operations use read locks for better performance
- Write operations use exclusive locks
- Tested with concurrent goroutines

## 🔮 Roadmap

- [ ] Connection pooling for clients
- [ ] TTL (Time To Live) support
- [ ] Persistence options
- [ ] Metrics and monitoring
- [ ] REST API interface
- [ ] Configuration file support
- [ ] Clustering support
- [ ] Memory usage optimization

### Development Setup
```bash
git clone https://github.com/ayushvyas-1/gcache.git
cd gcache
go mod tidy
make build
make test
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Redis and Memcached
- Built with Go's excellent concurrency primitives
- Thanks to the Go community for best practices

---

**Made with ❤️ in Go**