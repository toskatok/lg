/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-09-2018
 * |
 * | File Name:     generators/atrovan.go
 * +===============================================
 */

package generators

import (
	"encoding/json"
	"fmt"
)

// AtrovanGenerator generates data based on
// atrovan protocol
// and atrovan model.
type AtrovanGenerator struct {
}

// Topic returns atrovan mqtt topic
func (g AtrovanGenerator) Topic() string {
	return "v1/devices/me/telemetry"
}

// Generate generates atrovan message by converting input into telemetries json
// and using generator parameters.
func (g AtrovanGenerator) Generate(input interface{}) ([]byte, error) {
	// input into json
	states, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Input must be a map between strings and values")
	}

	message, err := json.Marshal(states)
	if err != nil {
		return nil, err
	}

	return message, nil
}
