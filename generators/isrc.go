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
	"time"

	"github.com/2tvenom/cbor"
	"github.com/I1820/link/protocols/lora"
)

// ISRCGenerator generates data based on
// [lora structure](https://github.com/aiotrc/uplink/blob/master/lora/messages.go) protocol
// and [cbor](http://cbor.io/) model.
// for historical reasons for refer to it as ISRC protocol
type ISRCGenerator struct {
	ApplicationID   string `mapstructure:"applicationID"`
	ApplicationName string `mapstructure:"applicationName"`
	DeviceName      string `mapstructure:"deviceName"`
	DevEUI          string `mapstructure:"devEUI"`
	GatewayMac      string `mapstructure:"gatewayMAC"`
}

// Topic returns lora mqtt topic
func (g ISRCGenerator) Topic() string {
	return fmt.Sprintf("application/%s/device/%s/rx", g.ApplicationID, g.DevEUI)
}

// Generate generates lora message by converting input into cbor and generator
// parameters.
func (g ISRCGenerator) Generate(input interface{}) ([]byte, error) {
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
