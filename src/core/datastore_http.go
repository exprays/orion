package core

import (
	"orion/src/data"
	"strings"
)

// HTTP-specific methods for the DataStore

func GetKeyCount() int {
	return data.Store.DBSize()
}

func GetKeys(pattern string, limit int) []string {
	// Get all data and filter by pattern
	allData := data.Store.GetAllData()

	var keys []string
	count := 0

	for key := range allData {
		if count >= limit {
			break
		}

		if matchPattern(key, pattern) {
			keys = append(keys, key)
			count++
		}
	}

	return keys
}

func GetKeyType(key string) string {
	// Check in different stores to determine type
	if data.Store.Exists(key) {
		return "string"
	}

	// Check if it's a set
	members := data.Store.SMembers(key)
	if len(members) > 0 || data.Store.SCard(key) >= 0 {
		return "set"
	}

	// Check if it's a hash
	if data.Store.HLen(key) >= 0 {
		return "hash"
	}

	return ""
}

func GetString(key string) *string {
	if value, exists := data.Store.Get(key); exists {
		return &value
	}
	return nil
}

func GetSetMembers(key string) []string {
	return data.Store.SMembers(key)
}

func GetHashFields(key string) map[string]string {
	// Since the data store doesn't have a direct method to get all hash fields,
	// we'll need to implement this based on the hash structure
	// For now, return empty map - this will need to be enhanced
	return make(map[string]string)
}

func GetTTL(key string) int {
	ttl := data.Store.TTL(key)
	return int(ttl)
}

// Simple pattern matching for keys (supports * wildcard)
func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}

	// Simple wildcard matching
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			prefix, suffix := parts[0], parts[1]
			return strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix)
		}
	}

	return key == pattern
}
