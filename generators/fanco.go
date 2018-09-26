/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-08-2018
 * |
 * | File Name:     fanco.go
 * +===============================================
 */

package generators

import (
	"encoding/json"
	"fmt"
	"time"
)

// FancoGenerator generates data with allthingstalk format and
// sends them with pure json.
type FancoGenerator struct {
	ThingID string
}

// Topic returns I1820 thing state topic.
// this topic sets thing state (with all of its assets) in I1820
func (g FancoGenerator) Topic() string {
	return fmt.Sprintf("things/%s/state", g.ThingID)
}

// Generate generates data message (in thing state format with all of its assets)
// in pure json.
func (g FancoGenerator) Generate(input interface{}) ([]byte, error) {
	// convert given input into json
	values, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Input must be a map between strings and values")
	}

	// allthingstalk state format
	// {
	//    "temperature": {
	//      "value": 10,
	//      "at": ...
	//    }
	// }
	states := make(map[string]struct {
		At    time.Time
		Value interface{}
	})
	for name, value := range values {
		states[name] = struct {
			At    time.Time
			Value interface{}
		}{time.Now(), value}
	}

	b, err := json.Marshal(states)
	if err != nil {
		return nil, err
	}

	return b, nil
}
