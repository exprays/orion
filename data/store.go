package data

import (
    "strconv"
    "sync"
    "time"
)

// DataStore represents the in-memory key-value store
type DataStore struct {
    mu    sync.RWMutex
    store map[string]string
    TTLStore map[string]int64 // Stores TTL (Time to Live) for each key in seconds
}

// Store is the global instance of DataStore
var Store *DataStore

func init() {
    Store = NewDataStore()
}

// NewDataStore initializes a new data store
func NewDataStore() *DataStore {
    return &DataStore{
        store: make(map[string]string),
        TTLStore: make(map[string]int64),
    }
    go ds.startExpirationCheck()
    return ds
}

// Implement a goroutine to periodically check and remove expired keys
func (ds *DataStore) startExpirationCheck() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        <-ticker.C
        ds.mu.Lock()
        for key, ttl := range ds.TTLStore {
            if ttl <= 0 {
                delete(ds.store, key)
                delete(ds.TTLStore, key)
            } else {
                ds.TTLStore[key] = ttl - 1
            }
        }
        ds.mu.Unlock()
    }
}

// Set stores a value associated with a key
func (ds *DataStore) Set(key, value string) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    ds.store[key] = value
}

// Get retrieves a value associated with a key
func (ds *DataStore) Get(key string) (string, bool) {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    value, exists := ds.store[key]
    return value, exists
}

// Append appends a value to the string at the specified key
func (ds *DataStore) Append(key, value string) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    if existing, exists := ds.store[key]; exists {
        ds.store[key] = existing + value
    } else {
        ds.store[key] = value
    }
}

// DecrBy decrements the integer value of a key by the given number
func (ds *DataStore) DecrBy(key string, decrement int) int {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    existing, exists := ds.store[key]
    if !exists {
        ds.store[key] = strconv.Itoa(-decrement)
        return -decrement
    }
    value, err := strconv.Atoi(existing)
    if err != nil {
        ds.store[key] = strconv.Itoa(-decrement)
        return -decrement
    }
    value -= decrement
    ds.store[key] = strconv.Itoa(value)
    return value
}

// Decr decrements the integer value of a key by 1
func (ds *DataStore) Decr(key string) int {
    return ds.DecrBy(key, 1)
}

// Del deletes a key from the store
func (ds *DataStore) Del(key string) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    delete(ds.store, key)
}

// GetDel retrieves a value associated with a key and deletes the key
func (ds *DataStore) GetDel(key string) (string, bool) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    value, exists := ds.store[key]
    if exists {
        delete(ds.store, key)
    }
    return value, exists
}

// GetEx retrieves a value associated with a key and sets expiration in seconds
func (ds *DataStore) GetEx(key string, seconds int64) (string, bool) {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    value, exists := ds.store[key]
    if exists && seconds > 0 {
        ds.TTLStore[key] = seconds
    }
    return value, exists
}

// GetRange retrieves a substring of the string value stored at a key
func (ds *DataStore) GetRange(key string, start, end int) string {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    value, exists := ds.store[key]
    if !exists {
        return ""
    }
    // Calculate proper start and end indices
    strLen := len(value)
    if start < 0 {
        start = strLen + start
    }
    if end < 0 {
        end = strLen + end
    }
    if start < 0 {
        start = 0
    }
    if end > strLen-1 {
        end = strLen - 1
    }
    if start > end {
        return ""
    }
    return value[start : end+1]
}


// FlushAll clears all key-value pairs from the store
func (ds *DataStore) FlushAll() {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    ds.store = make(map[string]string)
}
