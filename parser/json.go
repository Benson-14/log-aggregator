package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Parse(line string) (*LogEntry, error) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, err
	}

	entry := &LogEntry{
		Fields: make(map[string]string),
	}

	if t, ok := raw["time"].(string); ok {
		entry.Timestamp, _ = time.Parse(time.RFC3339, t)
	}

	if level, ok := raw["level"].(string); ok {
		entry.Level = level
	}

	if source, ok := raw["source"].(string); ok {
		entry.Source = source
	}

	if msg, ok := raw["message"].(string); ok {
		entry.Message = msg
	}

	for k, v := range raw {
		switch k {
		case "time", "level", "source", "message":
			continue
		default:
			entry.Fields[k] = fmt.Sprintf("%v", v)
		}
	}
	return entry, nil

}
