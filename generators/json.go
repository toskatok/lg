package generators

import (
	"encoding/json"
)

// JSONGenerator generates data with json format
type JSONGenerator struct {
	Name string `mapstructure:"name"`
}

// Topic returns given name as topic
func (g JSONGenerator) Topic() string {
	return g.Name
}

// Generate generates data message json
func (g JSONGenerator) Generate(input interface{}) ([]byte, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return b, nil
}
