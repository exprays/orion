
# Orion ![Orion Logo](https://img.shields.io/badge/Orion-Database-blue?style=for-the-badge&logo=database)

> An open-source, insanely-fast, in-memory database built with Go

![Go Version](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)
![Build Status](https://img.shields.io/github/actions/workflow/status/exprays/orion/build.yml?branch=main&style=for-the-badge)

---

Orion is a high-performance, in-memory key-value database designed for ultra-fast data access and manipulation. Whether you're building real-time applications, caching layers, or need rapid data processing, Orion is engineered to handle it with elegance and speed.

> ⚠️ **Note**: Orion is currently in active development as part of Buildspace S5. Expect frequent updates, rapid improvements, and exciting new features.

---

## 📖 Table of Contents
- [Features](#-features)
- [Quick Start](#-quick-start)
- [Usage Examples](#-usage-examples)
- [Benchmarks](#-benchmarks)
- [Architecture](#-architecture)
- [ORSP Protocol](#-orsp-protocol)
- [Tech Stack](#-tech-stack)

---

## ✨ Features

### 🔥 Core Features

| Feature              | Status | Description                              |
|----------------------|--------|------------------------------------------|
| In-memory KV Store   | ✅     | Lightning-fast key-value operations      |
| Strings              | ✅     | Full string manipulation support         |
| Sets                 | ✅     | Efficient set operations and manipulations |
| Persistence (AOF)    | ✅     | Append-only file for data durability     |
| TTL Support          | ✅     | Automatic key expiration                 |
| CLI Client (Hunter)  | ✅     | Interactive command-line interface       |
| Custom Protocol (ORSP)| ✅    | Optimized binary/text serialization      |
| Enhanced Logging     | ✅     | Error, Command logging for Server        |

### 🎯 Advanced Features

| Feature              | Status | Description                              |
|----------------------|--------|------------------------------------------|
| Command Autocomplete | ✅     | Tab completion with case-insensitive match |
| Command History      | ✅     | Persistent history navigation            |
| Background Saves     | ✅     | Non-blocking snapshots                   |
| AOF Rewriting        | ✅     | Log compaction with background safety    |
| Server Monitoring    | ❌    | Real-time statistics and metrics         |

### 🚧 Coming Soon

| Feature              | Status | ETA      | Priority |
|----------------------|--------|----------|----------|
| Sorted Sets          | 🔄     | Q1 2025  | High     |
| Hash Maps            | 🔄     | Q1 2025  | High     |
| Pub/Sub              | 📋     | Q2 2025  | Medium   |
| Transactions         | 📋     | Q2 2025  | Medium   |
| Streams              | 📋     | Q2 2025  | Medium   |
| Clustering           | 📋     | Q2 2025  | High     |
| Authentication       | 📋     | Q2 2025  | Medium   |
| LRU Eviction         | 📋     | Q2 2025  | Medium   |
| HyperLogLogs         | 📋     | Q3 2025  | Low      |
| Bitmaps              | 📋     | Q3 2025  | Low      |
| Vector Support       | 📋     | Q4 2025  | Low      |
| Geo-spatial Data     | 📋     | Q4 2025  | Low      |

---

## 🚀 Quick Start

### 🧰 Prerequisites
- Go 1.22 or higher installed
- Basic familiarity with terminal/CLI

### 🔧 Installation

#### Option 1: Clone from Source

```bash
git clone https://github.com/exprays/orion.git
cd orion
go mod tidy
```

#### Option 2: Use Prebuilt Binary

Download the latest binary from the [GitHub Releases](https://github.com/exprays/orion/releases).

---

### 🖥 Running Orion

#### Start the Server

```bash
cd cmd/server
go run orion.go --port=6379

# Or just go run orion.go which starts the server in default port 6379
```


#### Launch the Hunter CLI

```bash
cd cmd/hunter
go run hunter.go

# Select dev for default port 
# Select custom for custom port
```

> Press `CTRL+C` to exit the client.

---

## 🎮 Usage Examples

### Basic Operations

```bash
# String Operations

orion> SET user:1 "John Doe"
OK

orion> GET user:1
"John Doe"

orion> APPEND user:1 " - Engineer"
(integer) 19

# Set Operations

orion> SADD languages go python rust
(integer) 3

orion> SMEMBERS languages
1) "go"
2) "python"
3) "rust"

# TTL Operations

orion> SET session:abc123 "active" EX 3600
OK

orion> TTL session:abc123
(integer) 3599

```

### Advanced Features

```bash
# Command autocomplete (press TAB)
orion> s<TAB>
SADD  SCARD  SDIFF  SDIFFSTORE  SET  SISMEMBER  SMEMBERS  SMOVE  SPOP  SRANDMEMBER  SREM  SUNION  SUNIONSTORE

# Command history (use ↑/↓ arrows)
orion> <UP ARROW>  # Shows previous command

# Server information
orion> INFO
# Server
uptime_in_seconds:120
uptime_in_days:0
# Memory
used_memory:1048576
used_memory_human:1.0 MB
# Keyspace
db0:keys=5
```

### Autocomplete & History

- Use `<TAB>` for suggestions
- Use `↑ / ↓` arrows to navigate command history

---

## 📊 Benchmarks

| Operation | Ops/sec | Avg Latency | P99 Latency |
|----------:|--------:|-------------|-------------|
| SET       | 180,000 | 0.05ms      | 0.2ms       |
| GET       | 220,000 | 0.04ms      | 0.15ms      |
| SADD      | 150,000 | 0.06ms      | 0.25ms      |
| SMEMBERS  | 90,000  | 0.1ms       | 0.4ms       |

**Memory Usage**:
- Startup: ~8MB
- Per Key: ~64 bytes
- Per Set Member: ~32 bytes
- AOF File: ~40% of total data size

### Throughput Comparison

**SET Operation (10M)**

| Database | Ops/sec | Memory Usage |
|----------|---------|--------------|
| Orion    | 180,000 | 245MB        |
| Redis    | 165,000 | 280MB        |
| KeyDB    | 170,000 | 260MB        |

---

## 🧠 Architecture

```text
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Hunter CLI  │───▶ │ ORSP Protocol│───▶│ Orion Server │
│ (Client App) │     │ (Serializer) │     │ (Database)   │
└──────────────┘     └──────────────┘     └──────────────┘
                                               │
                                  ┌────────────────────┐
                                  │     Data Store     │
                                  │ ┌───────────────┐  │
                                  │ │ KV Store      │  |
                                  | | Set Store     |  |
                                  | | TTL Store     |  |
                                  │ └───────────────┘  │
                                  └────────────────────┘
                                               │
                                  ┌────────────────────┐
                                  │   Persistence      │
                                  │ ┌───────────────┐  │
                                  │ │ AOF Writer    │  │
                                  | | Snapshots     |  |
                                  │ └───────────────┘  │
                                  └────────────────────┘
```

---

# 📦 ORSP Protocol

## 🔄 ORSP Type Symbols and Mappings

 The ORSP (Orion Serialization Protocol) types and their corresponding Go type representations.


| ORSP Type Symbol | ORSP Name             | Go Type                | Description                                                                 |
|------------------|-----------------------|------------------------|-----------------------------------------------------------------------------|
| `+`              | Simple String         | `SimpleStringValue`    | Represents simple text response (e.g., `+OK\r\n`)                         |
| `-`              | Error                 | `ErrorValue`           | Represents an error message (e.g., `-ERR something went wrong\r\n`)      |
| `:`              | Integer               | `IntegerValue`         | Represents integer values (e.g., `:123\r\n`)                              |
| `$`              | Bulk String           | `BulkStringValue`      | Represents strings with length prefix (e.g., `$6\r\nfoobar\r\n`)       |
| `*`              | Array                 | `ArrayValue`           | Represents an array of ORSP values                                          |
| `_`              | Null                  | `NullValue`            | Represents a null value                                                     |
| `#`              | Boolean               | `BooleanValue`         | Represents a boolean (`#t\r\n` or `#f\r\n`)                             |
| `,`              | Double                | `DoubleValue`          | Represents a double-precision float                                         |
| `(`              | Big Number            | `BigNumberValue`       | Represents large integers using `math/big.Int`                             |
| `!`              | Bulk Error            | `BulkErrorValue`       | Contains structured error with code and message                            |
| `=`              | Verbatim String       | `VerbatimStringValue`  | Format-qualified string (e.g., `=txt:Hello World\r\n`)                    |
| `%`              | Map                   | `MapValue`             | Key-value mapping, keys must be simple strings                             |
| `~`              | Set                   | `SetValue`             | Unordered collection of unique elements                                     |
| `>`              | Push                  | `PushValue`            | Push-type message with kind and data array                                 |

---

Each type supports a `Marshal()` method to serialize the value to ORSP-compliant wire format, and a corresponding unmarshal function to deserialize from a stream.

                               |

---


## 💻 Tech Stack

- ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
- ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

---

## 📜 License
This project is licensed under the MIT License - see the LICENSE file for details.
MIT © [exprays](https://github.com/exprays/orion)

---

## 🙏 Acknowledgments
- Built with ❤️ by Surya A.K.A exprays
- Part of Buildspace S5 cohort
- Inspired by Redis and modern in-memory databases
- Special thanks to the Go community

## 📞 Support & Community
- Documentation: orion.thestarsociety.tech
- Issues: GitHub Issues
- Discussions: GitHub Discussions
- Instagram: @suryakantsubudhi
