package data

import (
	"fmt"
	"orion/src/aof"
	"orion/src/protocol"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// DataStore represents the in-memory key-value store
type DataStore struct {
	mu        sync.RWMutex
	store     map[string]string
	setStore  map[string]map[string]struct{} // Field for sets
	TTLStore  map[string]int64               // Stores TTL (Time to Live) for each key in seconds
	startTime time.Time
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
		setStore: make(map[string]map[string]struct{}), // Initialize setStore ( for implementation of sets )
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
func (ds *DataStore) Set(key, value string, ttl time.Duration) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.store[key] = value

	// Set or remove TTL
	if ttl > 0 {
		ds.TTLStore[key] = int64(ttl / time.Second)
	} else {
		delete(ds.TTLStore, key)
	}

	// Append to AOF
	// Prepare ORSP command for AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SET"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(value),
	}

	if ttl > 0 {
		command = append(command,
			protocol.BulkStringValue("PX"),
			protocol.BulkStringValue(fmt.Sprintf("%d", ttl.Milliseconds())),
		)
	}

	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}

// Exists checks if a key exists in the store
func (ds *DataStore) Exists(key string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	_, exists := ds.store[key]
	return exists
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
	// Prepare ORSP command for AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("APPEND"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(value),
	}

	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}

// DecrBy decrements the integer value of a key by the given number
func (ds *DataStore) DecrBy(key string, decrement int) (int, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	existing, exists := ds.store[key]
	var value int
	var err error

	if !exists {
		value = -decrement
	} else {
		value, err = strconv.Atoi(existing)
		if err != nil {
			return 0, fmt.Errorf("value is not an integer")
		}
		value -= decrement
	}

	ds.store[key] = strconv.Itoa(value)

	// Prepare ORSP command for AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("DECRBY"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(strconv.Itoa(decrement)),
	}

	if err := aof.AppendCommand(command); err != nil {
		return value, fmt.Errorf("error appending to AOF: %w", err)
	}

	return value, nil
}

// Decr decrements the integer value of a key by 1
func (ds *DataStore) Decr(key string) (int, error) {
	// Prepare ORSP command for AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("DECR"),
		protocol.BulkStringValue(key),
	}

	if err := aof.AppendCommand(command); err != nil {
		return 0, fmt.Errorf("error appending to AOF: %w", err)
	}

	return ds.DecrBy(key, 1)
}

// Del deletes a key from the store
func (ds *DataStore) Del(key string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.store, key)

	// Append to AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("DEL"),
		protocol.BulkStringValue(key),
	}

	if err := aof.AppendCommand(command); err != nil {
		fmt.Errorf("error appending to AOF: %w", err)
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
		command := protocol.ArrayValue{
			protocol.BulkStringValue("GETDEL"),
			protocol.BulkStringValue(key),
		}
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
		command := protocol.ArrayValue{
			protocol.BulkStringValue("GETEX"),
			protocol.BulkStringValue(key),
			protocol.BulkStringValue(seconds),
		}
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
	command := protocol.ArrayValue{
		protocol.BulkStringValue("GETSET"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(value),
	}
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
		command := protocol.ArrayValue{
			protocol.BulkStringValue("INCR"),
			protocol.BulkStringValue(key),
		}
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
	command := protocol.ArrayValue{
		protocol.BulkStringValue("INCR"),
		protocol.BulkStringValue(key),
	}
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
		command := protocol.ArrayValue{
			protocol.BulkStringValue("INCRBY"),
			protocol.BulkStringValue(key),
			protocol.BulkStringValue(strconv.Itoa(increment)),
		}
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
	command := protocol.ArrayValue{
		protocol.BulkStringValue("INCRBY"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(strconv.Itoa(increment)),
	}
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
		command := protocol.ArrayValue{
			protocol.BulkStringValue("INCRBYFLOAT"),
			protocol.BulkStringValue(key),
			protocol.BulkStringValue(strconv.FormatFloat(increment, 'f', -1, 64)),
		}
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
	command := protocol.ArrayValue{
		protocol.BulkStringValue("INCRBYFLOAT"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(strconv.FormatFloat(floatValue, 'f', -1, 64)),
	}
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
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SETEX"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(seconds),
		protocol.BulkStringValue(value),
	}
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

// Time returns the current server time in seconds and microseconds
func (ds *DataStore) Time() (string, error) {
	now := time.Now()
	seconds := now.Unix()
	microseconds := now.UnixNano() / int64(time.Microsecond)
	// Format the response
	response := fmt.Sprintf("[%d %d]", seconds, microseconds%1e6)
	return response, nil
}

// Info gathers various statistics about the server and formats them
func (ds *DataStore) Info() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Get memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Collect keyspace information
	keyspaceInfo := ds.getKeyspaceInfo()

	// Format the information
	info := fmt.Sprintf(
		"# Server\n"+
			"uptime_in_seconds:%d\n"+
			"uptime_in_days:%d\n"+
			"# Memory\n"+
			"used_memory:%d\n"+
			"used_memory_human:%s\n"+
			"total_allocated_memory:%d\n"+
			"# Keyspace\n"+
			"%s",
		ds.GetUptimeSeconds(),
		ds.GetUptimeDays(),
		memStats.Alloc,
		humanReadableBytes(memStats.Alloc),
		memStats.TotalAlloc,
		keyspaceInfo,
	)

	return info
}

// GetUptimeSeconds returns the uptime of the server in seconds
func (ds *DataStore) GetUptimeSeconds() int64 {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Calculate the diff	uptime := time.Since(ds.startTime).Seconds()e
	uptime := time.Since(ds.startTime).Seconds()

	return int64(uptime)
}

// getUptimeDays returns the uptime of the server in days
func (ds *DataStore) GetUptimeDays() int64 {
	return ds.GetUptimeSeconds() / (60 * 60 * 24)
}

// getKeyspaceInfo collects keyspace information
func (ds *DataStore) getKeyspaceInfo() string {
	numKeys := len(ds.store)
	return fmt.Sprintf("db0:keys=%d", numKeys)
}

// humanReadableBytes converts bytes to a human-readable string
func humanReadableBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// DBSize returns the number of keys in the data store
func (ds *DataStore) DBSize() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Count keys in the main store
	mainStoreSize := len(ds.store)

	// Count keys in the set store
	setStoreSize := len(ds.setStore)

	// Return the total number of keys
	return mainStoreSize + setStoreSize
}

// FLUSHALL UNIVERSAL .....

// FlushAll clears all key-value pairs from the store
func (ds *DataStore) FlushAll() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.store = make(map[string]string)
	ds.setStore = make(map[string]map[string]struct{})
	ds.TTLStore = make(map[string]int64)

	// Append to AOF
	command := protocol.ArrayValue{protocol.BulkStringValue("FLUSHALL")}
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}
}

// ENDS HERE

// SETS
// DATATYPE: SET
//IN-MEMORY STORE IMPLEMENTATION OF SETS IN ORION

// SAdd adds the specified members to the set stored at key
func (ds *DataStore) SAdd(key string, members ...string) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, exists := ds.setStore[key]; !exists {
		ds.setStore[key] = make(map[string]struct{})
	}

	added := 0
	for _, member := range members {
		if _, exists := ds.setStore[key][member]; !exists {
			ds.setStore[key][member] = struct{}{}
			added++
		}
	}

	// Append to AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SADD"),
		protocol.BulkStringValue(key),
	}
	for _, member := range members {
		command = append(command, protocol.BulkStringValue(member))
	}
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return added
}

// SMembers returns all the members of the set value stored at key
func (ds *DataStore) SMembers(key string) []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	set, exists := ds.setStore[key]
	if !exists {
		return []string{}
	}

	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}

	return members
}

// SIsMember returns if member is a member of the set stored at key
func (ds *DataStore) SIsMember(key, member string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	set, exists := ds.setStore[key]
	if !exists {
		return false
	}

	_, isMember := set[member]
	return isMember
}

// SCard returns the cardinality (number of elements) of the set stored at key
func (ds *DataStore) SCard(key string) int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if set, exists := ds.setStore[key]; exists {
		return len(set)
	}
	return 0
}

// SMove moves member from the set at source to the set at destination
func (ds *DataStore) SMove(source, destination, member string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	sourceSet, sourceExists := ds.setStore[source]
	if !sourceExists {
		return false
	}

	if _, exists := sourceSet[member]; !exists {
		return false
	}

	if _, destExists := ds.setStore[destination]; !destExists {
		ds.setStore[destination] = make(map[string]struct{})
	}

	delete(sourceSet, member)
	ds.setStore[destination][member] = struct{}{}

	// Append to AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SMOVE"),
		protocol.BulkStringValue(source),
		protocol.BulkStringValue(destination),
		protocol.BulkStringValue(member),
	}
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return true
}

// SPop removes and returns one or more random members from the set
func (ds *DataStore) SPop(key string, count int) []string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	set, exists := ds.setStore[key]
	if !exists || len(set) == 0 {
		return nil
	}

	if count > len(set) {
		count = len(set)
	}

	members := make([]string, 0, count)
	for member := range set {
		if len(members) < count {
			members = append(members, member)
			delete(set, member)
		} else {
			break
		}
	}

	// Append to AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SPOP"),
		protocol.BulkStringValue(key),
		protocol.BulkStringValue(fmt.Sprintf("%d", count)),
	}
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return members
}

// SRem removes one or more members from the set
func (ds *DataStore) SRem(key string, members ...string) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	set, exists := ds.setStore[key]
	if !exists {
		return 0
	}

	removed := 0
	for _, member := range members {
		if _, exists := set[member]; exists {
			delete(set, member)
			removed++
		}
	}

	// Append to AOF
	command := protocol.ArrayValue{
		protocol.BulkStringValue("SREM"),
		protocol.BulkStringValue(key),
	}
	for _, member := range members {
		command = append(command, protocol.BulkStringValue(member))
	}
	if err := aof.AppendCommand(command); err != nil {
		fmt.Println("Error appending to AOF:", err)
	}

	return removed
}

// SDiff returns the difference between the sets stored at the given keys
func (ds *DataStore) SDiff(keys ...string) []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if len(keys) == 0 {
		return []string{}
	}

	result := ds.setStore[keys[0]]

	for _, key := range keys[1:] {
		set, exists := ds.setStore[key]
		if exists {
			for member := range result {
				if _, ok := set[member]; ok {
					delete(result, member)
				}
			}
		}
	}

	members := make([]string, 0, len(result))
	for member := range result {
		members = append(members, member)
	}

	return members
}

// SDiffStore stores the difference between the sets stored at the given keys in the destination key
func (ds *DataStore) SDiffStore(destination string, keys ...string) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if len(keys) == 0 {
		return 0
	}

	result := ds.setStore[keys[0]]

	for _, key := range keys[1:] {
		set, exists := ds.setStore[key]
		if exists {
			for member := range result {
				if _, ok := set[member]; ok {
					delete(result, member)
				}
			}
		}
	}

	if _, exists := ds.setStore[destination]; !exists {
		ds.setStore[destination] = make(map[string]struct{})
	}

	for member := range result {
		ds.setStore[destination][member] = struct{}{}
	}

	return len(result)
}
