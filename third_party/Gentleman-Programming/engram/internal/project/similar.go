package project

import (
	"sort"
	"strings"
)

// ProjectMatch represents a project name that is similar to a query string.
type ProjectMatch struct {
	Name      string // The existing project name
	MatchType string // "case-insensitive", "substring", or "levenshtein"
	Distance  int    // Levenshtein distance (0 for case-insensitive and substring matches)
}

// FindSimilar finds projects similar to the given name from a list of existing
// project names. Similarity is determined by three criteria:
//
//  1. Case-insensitive exact match (different case, same letters)
//  2. Substring containment (query is a substring of candidate or vice-versa)
//  3. Levenshtein distance ≤ maxDistance
//
// Exact matches (identical strings) are always excluded.
//
// Results are ordered: case-insensitive matches first, then substring matches,
// then levenshtein matches sorted by distance ascending.
func FindSimilar(name string, existing []string, maxDistance int) []ProjectMatch {
	if maxDistance < 0 {
		maxDistance = 0
	}

	nameLower := strings.ToLower(strings.TrimSpace(name))

	// Scale maxDistance for short names to avoid noisy matches.
	// A 2-char name with maxDistance 3 would match almost everything.
	effectiveMax := maxDistance
	if len(nameLower) > 0 {
		halfLen := len(nameLower) / 2
		if halfLen < 1 {
			halfLen = 1
		}
		if effectiveMax > halfLen {
			effectiveMax = halfLen
		}
	}

	var caseMatches []ProjectMatch
	var subMatches []ProjectMatch
	var levMatches []ProjectMatch

	seen := make(map[string]bool)

	for _, candidate := range existing {
		// Skip exact match (same string, no drift)
		if candidate == name {
			continue
		}

		candidateLower := strings.ToLower(strings.TrimSpace(candidate))

		// Skip after case-fold too — that would be a normalized duplicate
		if candidateLower == nameLower {
			// Only add if the strings differ (different casing is still drift)
			if candidate != name {
				if !seen[candidate] {
					seen[candidate] = true
					caseMatches = append(caseMatches, ProjectMatch{
						Name:      candidate,
						MatchType: "case-insensitive",
						Distance:  0,
					})
				}
			}
			continue
		}

		// Substring match — skip for very short names (< 3 chars)
		// to avoid noisy matches like "go" matching "golang-tools"
		if len(nameLower) >= 3 {
			if strings.Contains(candidateLower, nameLower) || strings.Contains(nameLower, candidateLower) {
				if !seen[candidate] {
					seen[candidate] = true
					subMatches = append(subMatches, ProjectMatch{
						Name:      candidate,
						MatchType: "substring",
						Distance:  0,
					})
				}
				continue
			}
		}

		// Levenshtein distance (using scaled effectiveMax)
		dist := levenshtein(nameLower, candidateLower)
		if dist <= effectiveMax {
			if !seen[candidate] {
				seen[candidate] = true
				levMatches = append(levMatches, ProjectMatch{
					Name:      candidate,
					MatchType: "levenshtein",
					Distance:  dist,
				})
			}
		}
	}

	// Sort levenshtein results by distance ascending
	sort.Slice(levMatches, func(i, j int) bool {
		return levMatches[i].Distance < levMatches[j].Distance
	})

	result := make([]ProjectMatch, 0, len(caseMatches)+len(subMatches)+len(levMatches))
	result = append(result, caseMatches...)
	result = append(result, subMatches...)
	result = append(result, levMatches...)
	return result
}

// levenshtein computes the Levenshtein (edit) distance between strings a and b.
// Uses the standard dynamic-programming approach with O(min(|a|,|b|)) space
// by only keeping two rows of the DP table at a time.
func levenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	la, lb := len(ra), len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Ensure a is the shorter string for space optimisation
	if la > lb {
		ra, rb = rb, ra
		la, lb = lb, la
	}

	prev := make([]int, la+1)
	curr := make([]int, la+1)

	for i := 0; i <= la; i++ {
		prev[i] = i
	}

	for j := 1; j <= lb; j++ {
		curr[0] = j
		for i := 1; i <= la; i++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			del := prev[i] + 1
			ins := curr[i-1] + 1
			sub := prev[i-1] + cost
			curr[i] = minOf3(del, ins, sub)
		}
		prev, curr = curr, prev
	}

	return prev[la]
}

// minOf3 returns the smallest of three integers.
func minOf3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}
