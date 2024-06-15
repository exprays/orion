package data

import "sync"

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
