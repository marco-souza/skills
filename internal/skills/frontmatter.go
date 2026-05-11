package skills

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFrontmatter extracts YAML frontmatter from a markdown string.
// It returns the parsed frontmatter map, the remaining body content, and an error if parsing fails.
func ParseFrontmatter(content string) (map[string]any, string, error) {
	content = strings.TrimSpace(content)

	if !strings.HasPrefix(content, "---") {
		return nil, content, fmt.Errorf("missing YAML frontmatter (must start with ---)")
	}

	parts := strings.SplitN(content[3:], "---", 2)
	if len(parts) < 2 {
		return nil, content, fmt.Errorf("missing closing --- for frontmatter")
	}

	yamlContent := strings.TrimSpace(parts[0])
	body := strings.TrimSpace(parts[1])

	var fm map[string]any
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return nil, body, fmt.Errorf("invalid YAML: %w", err)
	}

	return fm, body, nil
}
