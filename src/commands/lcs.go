// commands/lcs.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strings"
)

// HandleLCS calculates the longest common subsequence between the value of a key and a given string
func HandleLCS(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'lcs' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	compareString, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid compare string")
	}

	value, exists := data.Store.Get(string(key))
	if !exists {
		return protocol.NullValue{}
	}

	lcs := longestCommonSubsequence(value, string(compareString))
	return protocol.BulkStringValue(lcs)
}

// longestCommonSubsequence computes the LCS of two strings
func longestCommonSubsequence(a, b string) string {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Reconstruct the LCS string
	var lcs strings.Builder
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			lcs.WriteByte(a[i-1])
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	// Reverse the LCS string because we constructed it backwards
	runes := []rune(lcs.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
