package json

import (
	"encoding/json"
)

// Generator generates data with json format
type Generator struct {
	Name string `mapstructure:"name"`
}

// Topic returns given name as topic
func (g Generator) Topic() string {
	return g.Name
}

// Generate generates data message json
func (g Generator) Generate(input interface{}) ([]byte, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return b, nil
}
