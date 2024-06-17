package frontmatter

import (
	"bufio"
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

func Parse(scanner *bufio.Scanner) (*Frontmatter, string, int, error) {
	if scanner.Scan() {
		line := scanner.Text()
		if line != "---" {
			return nil, line, 1, nil
		}
	} else {
		if err := scanner.Err(); err != nil {
			return nil, "", 1, fmt.Errorf("scan a template file: %w", err)
		}
		return nil, "", 1, nil
	}

	matterLines := make([]string, 0, 8) //nolint:mnd
	breaked := false
	ln := 1
	line := ""
	for scanner.Scan() {
		ln++
		line = scanner.Text()
		if breaked {
			if line == "" {
				continue
			}
			break
		}
		if line == "---" {
			breaked = true
			continue
		}
		matterLines = append(matterLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, line, ln, fmt.Errorf("scan a template file: %w", err)
	}

	matterS := strings.Join(matterLines, "\n")
	m := &Frontmatter{}
	if err := yaml.Unmarshal([]byte(matterS), m); err != nil {
		return nil, line, ln, fmt.Errorf("unmarshal a front matter as YAML: %w", err)
	}
	return m, line, ln, nil
}
