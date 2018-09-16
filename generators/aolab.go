/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-08-2018
 * |
 * | File Name:     aolab.go
 * +===============================================
 */

package generators

import (
	"encoding/json"
	"fmt"

	"github.com/I1820/lanserver/models"
	"github.com/I1820/link/models/aolab"
)

// AolabGenerator generates data based on
// lanserver protocol
// and aolab model.
type AolabGenerator struct {
	DevEUI string
}

// Topic returns lanserver mqtt topic
func (g AolabGenerator) Topic() []byte {
	return []byte(fmt.Sprintf("device/%s/rx", g.DevEUI))
}

// Generate generates lanserver message by converting input into json and using generator
// parameters.
func (g AolabGenerator) Generate(input interface{}) ([]byte, error) {
	// input into json
	states, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Input must be a map between strings and values")
	}
	b, err := json.Marshal(aolab.Log{
		Device: g.DevEUI,
		States: states,
	})
	if err != nil {
		return nil, err
	}

	// lanserver message
	message, err := json.Marshal(models.RxMessage{
		DevEUI: g.DevEUI,
		Data:   b,
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
