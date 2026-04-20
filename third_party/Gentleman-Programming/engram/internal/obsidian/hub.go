package obsidian

import (
	"fmt"
	"strings"
)

// ObsRef is a lightweight reference to an observation for use in hub notes.
// It carries only the fields needed to build wikilinks and type annotations.
type ObsRef struct {
	Slug     string // filename slug (without .md extension), e.g. "fixed-auth-bug-1"
	Title    string // human-readable title
	TopicKey string // observation's topic_key (may be empty)
	Type     string // observation type (e.g. "bugfix", "architecture")
}

// ShouldCreateTopicHub reports whether a topic prefix has enough observations
// to warrant creating a hub note. The threshold is ≥2 (REQ-EXPORT-05).
func ShouldCreateTopicHub(count int) bool {
	return count >= 2
}

// SessionHubMarkdown generates the markdown content for a session hub note.
// It lists all observations in the session as wikilinks.
//
// Output path: {vault}/engram/_sessions/{sessionID}.md
func SessionHubMarkdown(sessionID string, obs []ObsRef) string {
	var sb strings.Builder

	// ── YAML Frontmatter ──────────────────────────────────────────────────────
	sb.WriteString("---\n")
	sb.WriteString("type: session-hub\n")
	fmt.Fprintf(&sb, "session_id: %s\n", sessionID)
	sb.WriteString("tags:\n  - session\n")
	sb.WriteString("---\n")

	// ── Title ─────────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "# Session: %s\n", sessionID)

	// ── Observations list ─────────────────────────────────────────────────────
	sb.WriteString("\n## Observations\n")
	for _, ref := range obs {
		fmt.Fprintf(&sb, "- [[%s]]\n", ref.Slug)
	}

	return sb.String()
}

// TopicHubMarkdown generates the markdown content for a topic cluster hub note.
// It lists all observations sharing the same topic prefix as wikilinks with
// type annotations.
//
// Output path: {vault}/engram/_topics/{prefix}.md
// where prefix uses "--" instead of "/" for filesystem safety.
func TopicHubMarkdown(prefix string, obs []ObsRef) string {
	var sb strings.Builder

	// ── YAML Frontmatter ──────────────────────────────────────────────────────
	sb.WriteString("---\n")
	sb.WriteString("type: topic-hub\n")
	fmt.Fprintf(&sb, "topic_prefix: %s\n", prefix)
	sb.WriteString("tags:\n  - topic\n")
	sb.WriteString("---\n")

	// ── Title ─────────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "# Topic: %s\n", prefix)

	// ── Related Observations list ─────────────────────────────────────────────
	sb.WriteString("\n## Related Observations\n")
	for _, ref := range obs {
		fmt.Fprintf(&sb, "- [[%s]] (%s)\n", ref.Slug, ref.Type)
	}

	return sb.String()
}
