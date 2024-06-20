package data

import (
	"fmt"
	"orion/src/aof"
	"strconv"
	"sync"
	"time"
)

// DataStore represents the in-memory key-value store
type DataStore struct {
	mu       sync.RWMutex
	store    map[string]string
	TTLStore map[string]int64 // Stores TTL (Time to Live) for each key in seconds
}

// Store is the global instance of DataStore
var Store *DataStore

func init() {
	Store = NewDataStore()
	go Store.startExpirationCheck() // Start the expiration check goroutine for Store
}

// NewDataStore initializes a new data store
func NewDataStore() *DataStore {
	ds := &DataStore{
		store:    make(map[string]string),
		TTLStore: make(map[string]int64),
	}
	return ds
}

// startExpirationCheck is a goroutine to periodically check and remove expired keys
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

	// Append to AOF
	command := fmt.Sprintf("SET %s %s", key, value)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
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

	// Append to AOF
	command := fmt.Sprintf("APPEND %s %s", key, value)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
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

	// Append to AOF
	command := fmt.Sprintf("DECRBY %s %d", key, decrement)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
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

	// Append to AOF
	command := fmt.Sprintf("DEL %s", key)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}

// GetDel retrieves a value associated with a key and deletes the key
func (ds *DataStore) GetDel(key string) (string, bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	value, exists := ds.store[key]
	if exists {
		delete(ds.store, key)

		// Append to AOF
		command := fmt.Sprintf("GETDEL %s", key)
		if err := aof.AppendCommand(command); err != nil {
			fmt.Println("Error appending to AOF:", err)
		}
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

		// Append to AOF
		command := fmt.Sprintf("GETEX %s %d", key, seconds)
		if err := aof.AppendCommand(command); err != nil {
			fmt.Println("Error appending to AOF:", err)
		}
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

// GetSet sets a new value for a key and returns its old value
func (ds *DataStore) GetSet(key, value string) (string, bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	oldValue, exists := ds.store[key]
	ds.store[key] = value

	// Append to AOF
	command := fmt.Sprintf("GETSET %s %s", key, value)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return oldValue, exists
}

// Incr increments the integer value of a key by 1
func (ds *DataStore) Incr(key string) (int, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	value, exists := ds.store[key]
	if !exists {
		ds.store[key] = "1"
		// Append to AOF
		command := fmt.Sprintf("INCR %s", key)
		if err := aof.AppendCommand(command); err != nil {
			fmt.Println("Error appending to AOF:", err)
		}
		return 1, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value for key %s is not an integer", key)
	}

	intValue++
	ds.store[key] = strconv.Itoa(intValue)
	// Append to AOF
	command := fmt.Sprintf("INCR %s", key)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return intValue, nil
}

// IncrBy increments the integer value of a key by a specified amount
func (ds *DataStore) IncrBy(key string, increment int) (int, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	value, exists := ds.store[key]
	if !exists {
		ds.store[key] = strconv.Itoa(increment)
		// Append to AOF
		command := fmt.Sprintf("INCRBY %s %d", key, increment)
		if err := aof.AppendCommand(command); err != nil {
			fmt.Println("Error appending to AOF:", err)
		}
		return increment, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value for key %s is not an integer", key)
	}

	intValue += increment
	ds.store[key] = strconv.Itoa(intValue)
	// Append to AOF
	command := fmt.Sprintf("INCRBY %s %d", key, increment)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return intValue, nil
}

// IncrByFloat increments the float value of a key by a specified amount
func (ds *DataStore) IncrByFloat(key string, increment float64) (float64, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	value, exists := ds.store[key]
	if !exists {
		ds.store[key] = strconv.FormatFloat(increment, 'f', -1, 64)
		// Append to AOF
		command := fmt.Sprintf("INCRBYFLOAT %s %f", key, increment)
		if err := aof.AppendCommand(command); err != nil {
			fmt.Println("Error appending to AOF:", err)
		}
		return increment, nil
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("value for key %s is not a float", key)
	}

	floatValue += increment
	ds.store[key] = strconv.FormatFloat(floatValue, 'f', -1, 64)
	// Append to AOF
	command := fmt.Sprintf("INCRBYFLOAT %s %f", key, increment)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return floatValue, nil
}

// SetEx stores a value with a specified TTL (in seconds)
func (ds *DataStore) SetEx(key, value string, seconds int64) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.store[key] = value
	ds.TTLStore[key] = seconds

	// Append to AOF
	command := fmt.Sprintf("SETEX %s %d %s", key, seconds, value)
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}

// TTL retrieves the TTL of a key in seconds
func (ds *DataStore) TTL(key string) int64 {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.store[key]; !exists {
		return -1 // Key does not exist
	}
	ttl, exists := ds.TTLStore[key]
	if !exists {
		return -1 // Key exists but has no TTL
	}
	return ttl
}

// FlushAll clears all key-value pairs from the store
func (ds *DataStore) FlushAll() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.store = make(map[string]string)

	// Append to AOF
	command := "FLUSHALL"
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}
