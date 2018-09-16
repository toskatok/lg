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

// AolabGenerator generates data based on
// lanserver protocol
// and aolab model.
type AtrovanGenerator struct {
}

// Topic returns lanserver mqtt topic
func (g AtrovanGenerator) Topic() []byte {
	return []byte("v1/devices/me/telemetry")
}

// Generate generates lanserver message by converting input into cbor and generator
// parameters.
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
