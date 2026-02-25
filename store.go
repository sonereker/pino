package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

const pinoDir = ".pino"

func pinoDirPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, pinoDir)
}

func ensureDir() error {
	return os.MkdirAll(pinoDirPath(), 0755)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	prev := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prev = false
		} else if !prev {
			b.WriteRune('-')
			prev = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func savePino(summary, prompt, plan string) (string, error) {
	if err := ensureDir(); err != nil {
		return "", err
	}

	now := time.Now()
	slug := slugify(summary)
	if len(slug) > 50 {
		slug = slug[:50]
	}
	filename := fmt.Sprintf("%s_%s", now.Format("2006-01-02"), slug)

	// Avoid collisions
	base := filename
	for i := 1; ; i++ {
		path := filepath.Join(pinoDirPath(), filename+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}
		filename = fmt.Sprintf("%s-%d", base, i)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", summary)
	fmt.Fprintf(&b, "> %s\n\n", now.Format("2006-01-02 15:04"))
	b.WriteString("## Prompt\n\n")
	b.WriteString(prompt)
	b.WriteString("\n")
	if plan != "" {
		b.WriteString("\n## Plan\n\n")
		b.WriteString(plan)
		b.WriteString("\n")
	}

	path := filepath.Join(pinoDirPath(), filename+".md")
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		return "", err
	}
	return filename, nil
}

var dateLineRe = regexp.MustCompile(`^>\s*(.+)$`)

func parsePino(path string) (Pino, error) {
	f, err := os.Open(path)
	if err != nil {
		return Pino{}, err
	}
	defer f.Close()

	var p Pino
	p.Filename = strings.TrimSuffix(filepath.Base(path), ".md")

	scanner := bufio.NewScanner(f)
	var section string
	var promptLines, planLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# ") && p.Summary == "" {
			p.Summary = strings.TrimPrefix(line, "# ")
			continue
		}

		if m := dateLineRe.FindStringSubmatch(line); m != nil && p.CreatedAt.IsZero() {
			if t, err := time.Parse("2006-01-02 15:04", strings.TrimSpace(m[1])); err == nil {
				p.CreatedAt = t
			}
			continue
		}

		if line == "## Prompt" {
			section = "prompt"
			continue
		}
		if line == "## Plan" {
			section = "plan"
			continue
		}

		switch section {
		case "prompt":
			promptLines = append(promptLines, line)
		case "plan":
			planLines = append(planLines, line)
		}
	}

	p.Prompt = strings.TrimSpace(strings.Join(promptLines, "\n"))
	p.Plan = strings.TrimSpace(strings.Join(planLines, "\n"))

	if p.CreatedAt.IsZero() {
		info, _ := os.Stat(path)
		if info != nil {
			p.CreatedAt = info.ModTime()
		}
	}

	return p, scanner.Err()
}

func listAllPinos() ([]Pino, error) {
	entries, err := filepath.Glob(filepath.Join(pinoDirPath(), "*.md"))
	if err != nil {
		return nil, err
	}

	var pinos []Pino
	for _, path := range entries {
		p, err := parsePino(path)
		if err != nil {
			continue
		}
		pinos = append(pinos, p)
	}

	sort.Slice(pinos, func(i, j int) bool {
		return pinos[i].CreatedAt.After(pinos[j].CreatedAt)
	})

	if pinos == nil {
		pinos = []Pino{}
	}
	return pinos, nil
}

func searchPinos(keyword string) ([]Pino, error) {
	all, err := listAllPinos()
	if err != nil {
		return nil, err
	}

	keyword = strings.ToLower(keyword)
	results := []Pino{}
	for _, p := range all {
		haystack := strings.ToLower(p.Summary + " " + p.Prompt + " " + p.Plan)
		if strings.Contains(haystack, keyword) {
			results = append(results, p)
		}
	}
	return results, nil
}

func deletePino(filename string) error {
	path := filepath.Join(pinoDirPath(), filename+".md")
	return os.Remove(path)
}
