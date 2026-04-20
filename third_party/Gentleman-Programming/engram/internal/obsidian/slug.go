// Package obsidian implements the Engram → Obsidian vault export engine.
// It reads observations from the store and writes structured markdown files
// with YAML frontmatter, wikilinks, and hub notes.
package obsidian

import (
	"fmt"
	"regexp"
	"strings"
)

const maxSlugLen = 60

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts an observation title and ID into a filesystem-safe slug.
// The algorithm:
//  1. Lowercase the title
//  2. Replace non-alphanumeric characters with hyphens
//  3. Trim leading/trailing hyphens
//  4. Truncate to 60 chars (trimming trailing hyphens after truncation)
//  5. Append the ID for collision safety
//
// If the title is empty, the slug is "observation-{id}".
func Slugify(title string, id int64) string {
	if title == "" {
		return fmt.Sprintf("observation-%d", id)
	}

	s := strings.ToLower(title)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if len(s) > maxSlugLen {
		s = s[:maxSlugLen]
		s = strings.TrimRight(s, "-")
	}

	return fmt.Sprintf("%s-%d", s, id)
}
