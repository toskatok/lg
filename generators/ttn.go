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
)

// TTNGenerator generates data based on TheThingsNetwork format
// for I1820 ttn integration mdoule.
type TTNGenerator struct {
	ApplicationName string `mapstructure:"applicationName"`
	ApplicationID   string `mapstructure:"applicationID"`
	DevEUI          string `mapstructure:"devEUI"`
	DeviceName      string `mapstructure:"deviceName"`
}

// TTNRequest is a request structure of TTN integration module
type TTNRequest struct {
	AppID          string `json:"app_id"`
	DevID          string `json:"dev_id"`
	HardwareSerial string `json:"hardware_serial"`
	Port           int    `json:"port"`
	Counter        int    `json:"counter"`
	IsRetry        bool   `json:"is_retry"`
	Confirmed      bool   `json:"confirmed"`
	PayloadRaw     []byte `json:"payload_raw"`
	Metadata       struct {
		Time        time.Time `json:"time"`
		Frequency   float64   `json:"frequency"`
		Modulation  string    `json:"modulation"`
		DataRate    string    `json:"data_rate"`
		BitRate     int       `json:"bit_rate"`
		CondingRate string    `json:"conding_rate"`
		Gateways    []struct {
			GatewayID string `json:"gtw_id"`
		} `json:"gateways"`
	} `json:"metadata"`
}

// Topic returns ttn integration http topic
func (g TTNGenerator) Topic() string {
	return fmt.Sprintf("ttn/%s", g.ApplicationID)
}

// Generate generates lora message by converting input into cbor and generator
// parameters into ttn message format.
func (g TTNGenerator) Generate(input interface{}) ([]byte, error) {
	// input to cbor conversion
	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	if ok, err := encoder.Marshal(input); !ok {
		return nil, err
	}

	// ttn integration message + time
	request := TTNRequest{
		AppID:          g.ApplicationName,
		DevID:          g.DeviceName,
		HardwareSerial: g.DevEUI,
		PayloadRaw:     buffer.Bytes(),
	}
	request.Metadata.Time = time.Now()
	message, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return message, nil
}
