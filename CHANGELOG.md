# Changelog

All notable changes to this project will be documented in this file.

---

## [Unreleased]

### ✨ CLI Enhancements

- **Command Autocomplete**
  - Tab completion for all CLI commands (e.g., `SET`, `GET`, `SADD`)
  - Case-insensitive and fuzzy matching support
  - Works with both lowercase and uppercase variants

- **Command History**
  - Persistent command history stored at `~/.orion_history.json`
  - Navigation using arrow keys (↑/↓)
  - History retained across sessions
  - JSON format for easier inspection and debugging

---

## [v0.1.0] - Initial Release

### 🎯 Core Features

- **In-Memory Key-Value Store**
  - Minimalist KV server built from scratch
  - Supports basic in-memory operations with no external dependencies

- **Client Interface**
  - Basic CLI client to interact with the server
  - Input parsing and request dispatching

- **Command Handler**
  - Central routing logic for commands
  - Modular support for adding new commands

- **Custom Serialization Protocol (ORSP)**
  - Built from scratch serialization/deserialization logic in `orsp.go`
  - Supports multiple data types: Simple Strings, Bulk Strings, Integers, Arrays, Nulls, Errors, Sets, Maps, and more
  - Marshal and Unmarshal logic implemented for all ORSP types

- **Append-Only File (AOF) Persistence**
  - All executed commands are logged to `appendonly.orion`
  - On server restart, the AOF file is replayed to restore state
  - Duplicate command prevention using hash map
  - Graceful corruption recovery during command replay
  - AOF rewriting logic to compact file using current state snapshot

---

## 🔮 Planned Enhancements

- [ ] Command argument autocomplete  
- [ ] Ctrl+R history search  
- [ ] History pruning / max size limits  
- [ ] Customizable history file path  

---

## 🛡 Breaking Changes

- **None** — All new features are backward compatible.

---

## ⚙️ Performance Notes

- Minimal memory footprint for autocomplete and history
- AOF file is written on command execution and synced immediately
- Autocomplete performs instantly with the current command set
