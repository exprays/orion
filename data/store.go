package data

import (
    "strconv"
    "sync"
)

// DataStore represents the in-memory key-value store
type DataStore struct {
    mu    sync.RWMutex
    store map[string]string
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


// FlushAll clears all key-value pairs from the store
func (ds *DataStore) FlushAll() {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    ds.store = make(map[string]string)
}
