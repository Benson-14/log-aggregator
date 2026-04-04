package query

import "strings"

type Query struct {
	Fields map[string]string
	Text   string
}

func Parse(input string) *Query {
	q := &Query{
		Fields: make(map[string]string),
	}

	var textParts []string
	for _, part := range strings.Fields(input) {
		if idx := strings.Index(part, ":"); idx != -1 {
			key := part[:idx]
			value := part[idx+1:]
			q.Fields[key] = value
		} else {
			textParts = append(textParts, part)
		}
	}

	q.Text = strings.Join(textParts, " ")
	return q
}
