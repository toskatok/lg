/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-07-2018
 * |
 * | File Name:     lora.go
 * +===============================================
 */

package generators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/2tvenom/cbor"
	"github.com/aiotrc/uplink/lora"
)

// LoRaApplicationGenerator generates data based on
// [lora structure](https://github.com/aiotrc/uplink/blob/master/lora/messages.go)
// and [cbor](http://cbor.io/)
type LoRaApplicationGenerator struct {
	ApplicationID   string
	ApplicationName string
	DeviceName      string
	DevEUI          string
	GatewayMac      string
}

// Generate generates lora message by converting input into cbor and generator
// parameters.
func (g LoRaApplicationGenerator) Generate(input interface{}) ([]byte, error) {
	// input into cbor
	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	if ok, err := encoder.Marshal(input); !ok {
		return nil, err
	}

	// lora message
	message, err := json.Marshal(lora.RxMessage{
		ApplicationID:   g.ApplicationID,
		ApplicationName: g.ApplicationName,
		DeviceName:      g.DeviceName,
		DevEUI:          g.DevEUI,
		FPort:           5,
		FCnt:            10,
		RxInfo: []lora.RxInfo{
			lora.RxInfo{
				Mac:     g.GatewayMac,
				Name:    fmt.Sprintf("gateway-%s", g.GatewayMac),
				Time:    time.Now(),
				RSSI:    -57,
				LoRaSNR: 10,
			},
		},
		TxInfo: lora.TxInfo{
			Frequency: 868100000,
			Adr:       true,
			CodeRate:  "4/6",
		},
		Data: buffer.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
