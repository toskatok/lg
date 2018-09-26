/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-07-2018
 * |
 * | File Name:     isrc.go
 * +===============================================
 */

package generators

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/2tvenom/cbor"
	"github.com/I1820/link/protocols/lora"
)

// TTNGenerator generates data based on TheThingsNetwork format
// for I1820 ttn integration mdoule.
type TTNGenerator struct {
	ApplicationName string `mapstructure:"applicationName"`
	ApplicationID   string `mapstructure:"applicationID"`
	DevEUI          string `mapstructure:"devEUI"`
	DeviceName      string `mapstructure:"deviceName"`
}

// Topic returns ttn integration http topic
func (g TTNGenerator) Topic() string {
	return fmt.Sprintf("ttn/%s", g.ApplicationID)
}

// Generate generates lora message by converting input into cbor and generator
// parameters into ttn message format.
func (g TTNGenerator) Generate(input interface{}) ([]byte, error) {
	// input into cbor
	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	if ok, err := encoder.Marshal(input); !ok {
		return nil, err
	}

	// lora message
	message, err := json.Marshal(lora.RxMessage{
		ApplicationID: g.ApplicationName,
		DeviceName:    g.DeviceName,
		DevEUI:        g.DevEUI,
		FPort:         5,
		FCnt:          10,
		Data:          buffer.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
