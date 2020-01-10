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

package isrc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/2tvenom/cbor"
)

const (
	// FPort distinguishes the data format
	FPort = 5
	// FCnt counts number of frames
	FCnt = 10
	// Frequency of communication (Gateway)
	Frequency = 868100000
	// LoRaSNR indicates Signal to noise ratio (Gateway)
	LoRaSNR = 10
)

// RxMessage contains payloads received from your nodes in loraserver.io
type RxMessage struct {
	ApplicationID   string
	ApplicationName string
	DeviceName      string
	DevEUI          string
	FPort           int
	FCnt            int
	RxInfo          []RxInfo
	TxInfo          TxInfo
	Data            []byte
}

// RxInfo contains reception information of a lara gateway that
// payload is received from it.
type RxInfo struct {
	Mac     string
	Name    string
	Time    time.Time
	RSSI    int     `json:"rssi"`
	LoRaSNR float64 `json:"LoRaSNR"`
}

// TxInfo contains transmission information of a lora gateway that
// payload is received from it.
type TxInfo struct {
	Frequency int
	Adr       bool
	CodeRate  string
}

// Generator generates data based on
// RxMessage structure as is described above and then encode it
// with [cbor](http://cbor.io/).
// for historical reasons for refer to it as ISRC protocol
type Generator struct {
	ApplicationID   int    `mapstructure:"applicationID"`
	ApplicationName string `mapstructure:"applicationName"`
	DeviceName      string `mapstructure:"deviceName"`
	DevEUI          string `mapstructure:"devEUI"`
	GatewayMac      string `mapstructure:"gatewayMAC"`
}

// Topic returns lora mqtt topic
func (g Generator) Topic() string {
	return fmt.Sprintf("application/%d/device/%s/rx", g.ApplicationID, g.DevEUI)
}

// Generate generates lora message by converting input into cbor and generator
// parameters.
func (g Generator) Generate(input interface{}) ([]byte, error) {
	// input into cbor
	var buffer bytes.Buffer

	encoder := cbor.NewEncoder(&buffer)
	if ok, err := encoder.Marshal(input); !ok {
		return nil, err
	}

	// lora message
	message, err := json.Marshal(RxMessage{
		ApplicationID:   strconv.Itoa(g.ApplicationID),
		ApplicationName: g.ApplicationName,
		DeviceName:      g.DeviceName,
		DevEUI:          g.DevEUI,
		FPort:           FPort,
		FCnt:            FCnt,
		RxInfo: []RxInfo{
			{
				Mac:     g.GatewayMac,
				Name:    fmt.Sprintf("gateway-%s", g.GatewayMac),
				Time:    time.Now(),
				RSSI:    -57,
				LoRaSNR: LoRaSNR,
			},
		},
		TxInfo: TxInfo{
			Frequency: Frequency,
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
