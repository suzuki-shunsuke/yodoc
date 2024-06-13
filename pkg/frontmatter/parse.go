package frontmatter

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Dest  string
	Dir   string
	Env   map[string]string
	Delim *Delim
}

type Delim struct {
	Left  string
	Right string
}

func (d *Delim) GetLeft() string {
	if d == nil {
		return ""
	}
	return d.Left
}

func (d *Delim) GetRight() string {
	if d == nil {
		return ""
	}
	return d.Right
}

func Parse(s string) (*Frontmatter, string, error) {
	if !strings.HasPrefix(s, "---") {
		return nil, s, nil
	}
	lines := strings.Split(s, "\n")
	matterLines := make([]string, 0, 8) //nolint:mnd
	remaining := ""
	breaked := false
	for i, line := range lines[1:] {
		if breaked {
			if line == "" {
				continue
			}
			remaining = strings.Join(lines[i+1:], "\n")
			break
		}
		if line == "---" {
			breaked = true
			continue
		}
		matterLines = append(matterLines, line)
	}
	matterS := strings.Join(matterLines, "\n")
	m := &Frontmatter{}
	if err := yaml.Unmarshal([]byte(matterS), m); err != nil {
		return nil, "", fmt.Errorf("unmarshal a front matter as YAML: %w", err)
	}
	return m, remaining, nil
}
