package repository

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed queries/*.sql
var queryFiles embed.FS

// QueryLoader handles loading and parsing SQL queries from files
type QueryLoader struct {
	queries map[string]string
}

// NewQueryLoader creates a new query loader
func NewQueryLoader() (*QueryLoader, error) {
	loader := &QueryLoader{
		queries: make(map[string]string),
	}

	if err := loader.loadQueries(); err != nil {
		return nil, fmt.Errorf("loading queries: %w", err)
	}

	return loader, nil
}

// GetQuery returns the SQL query by name
func (ql *QueryLoader) GetQuery(name string) (string, error) {
	query, exists := ql.queries[name]
	if !exists {
		return "", fmt.Errorf("query %s not found", name)
	}
	return query, nil
}

func (ql *QueryLoader) loadQueries() error {
	files, err := queryFiles.ReadDir("queries")
	if err != nil {
		return fmt.Errorf("reading queries directory: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		content, err := queryFiles.ReadFile("queries/" + file.Name())
		if err != nil {
			return fmt.Errorf("reading query file %s: %w", file.Name(), err)
		}

		if err := ql.parseQueriesFromContent(string(content)); err != nil {
			return fmt.Errorf("parsing queries from %s: %w", file.Name(), err)
		}
	}

	return nil
}

func (ql *QueryLoader) parseQueriesFromContent(content string) error {
	lines := strings.Split(content, "\n")
	var currentQuery strings.Builder
	var currentName string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines/comments
		if line == "" || (strings.HasPrefix(line, "--") && !strings.Contains(line, "name:")) {
			continue
		}

		// Check if this is a query name definition
		if strings.HasPrefix(line, "-- name:") {
			if currentName != "" && currentQuery.Len() > 0 {
				ql.queries[currentName] = strings.TrimSpace(currentQuery.String())
			}

			parts := strings.SplitN(line, "name:", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid query name format: %s", line)
			}

			nameAndType := strings.TrimSpace(parts[1])
			nameParts := strings.Fields(nameAndType)
			if len(nameParts) < 1 {
				return fmt.Errorf("missing query name: %s", line)
			}

			currentName = nameParts[0]
			currentQuery.Reset()
			continue
		}

		if currentName != "" {
			if currentQuery.Len() > 0 {
				currentQuery.WriteString(" ")
			}
			currentQuery.WriteString(line)
		}
	}

	if currentName != "" && currentQuery.Len() > 0 {
		ql.queries[currentName] = strings.TrimSpace(currentQuery.String())
	}

	return nil
}
