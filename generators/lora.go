/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-01-2019
 * |
 * | File Name:     lora.go
 * +===============================================
 */

package generators

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/2tvenom/cbor"
	"github.com/brocaar/lorawan"
)

// LoRaGenerator generates data based on LoRaWAN protocol. It encrypts data and you will need
// a lora server to decode it.
type LoRaGenerator struct {
	Gateway struct {
		Mac string `mapstructure:"MAC"`
	} `mapstructure:"gateway"`
	Keys struct {
		NetworkSKey     string `mapstructure:"networkSKey"`
		ApplicationSKey string `mapstructure:"applicationSKey"`
	} `mapstructure:"keys"`
	Device struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"device"`
}

// DataRate contains information that gateway collects about packet data rate.
type DataRate struct {
	Modulation   string
	Bandwidth    int
	SpreadFactor int
}

// RxRawInfo is an information that is coming from the gateway
type RxRawInfo struct {
	Board     int
	Antenna   int
	Channel   int
	CodeRate  string
	CrcStatus int
	DataRate  DataRate
	Frequency int
	LoRaSNR   int
	Mac       string
	RfChain   int
	Rssi      int
	Size      int
}

// RxPacket is a packet that is sent to the loraserver.io based on
// https://www.loraserver.io/lora-gateway-bridge/use/data/
type RxPacket struct {
	PhyPayload []byte
	RxInfo     RxRawInfo
}

// Topic returns lora gateway mqtt topic.
func (g LoRaGenerator) Topic() string {
	return fmt.Sprintf("gateway/%s/rx", g.Gateway.Mac)
}

// Generate generates lora message by converting input into cbor and encrypts it.
func (g LoRaGenerator) Generate(input interface{}) ([]byte, error) {
	// encodes input with cbor
	var buffer bytes.Buffer

	encoder := cbor.NewEncoder(&buffer)
	if ok, err := encoder.Marshal(input); !ok {
		return nil, err
	}

	// converts network and application session keys to AES128
	appSKeySlice, err := hex.DecodeString(g.Keys.ApplicationSKey)
	if err != nil {
		return nil, err
	}

	var appSKey lorawan.AES128Key

	copy(appSKey[:], appSKeySlice)

	nwkSKeySlice, err := hex.DecodeString(g.Keys.NetworkSKey)
	if err != nil {
		return nil, err
	}

	var nwkSKey lorawan.AES128Key

	copy(nwkSKey[:], nwkSKeySlice)

	// converts device addr into DevAddr
	devAddrSlice, err := hex.DecodeString(g.Device.Addr)
	if err != nil {
		return nil, err
	}

	var devAddr lorawan.DevAddr

	copy(devAddr[:], devAddrSlice)

	// https://godoc.org/github.com/brocaar/lorawan#example-PHYPayload--Lorawan10Encode
	fport := uint8(5)

	phy := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.UnconfirmedDataUp,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: &lorawan.MACPayload{
			FHDR: lorawan.FHDR{
				DevAddr: devAddr,
				FCtrl: lorawan.FCtrl{
					ADR:       false,
					ADRACKReq: false,
					ACK:       false,
				},
				FCnt: 10,
				FOpts: []lorawan.Payload{
					&lorawan.MACCommand{
						CID: lorawan.DevStatusAns,
						Payload: &lorawan.DevStatusAnsPayload{
							Battery: 115,
							Margin:  7,
						},
					},
				},
			},
			FPort:      &fport,
			FRMPayload: []lorawan.Payload{&lorawan.DataPayload{Bytes: buffer.Bytes()}},
		},
	}
	if err := phy.EncryptFRMPayload(appSKey); err != nil {
		return nil, err
	}

	if err := phy.SetUplinkDataMIC(lorawan.LoRaWAN1_0, 0, 0, 0, nwkSKey, lorawan.AES128Key{}); err != nil {
		return nil, err
	}

	phyBytes, err := phy.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// lora message
	message, err := json.Marshal(RxPacket{
		RxInfo: RxRawInfo{
			Board:     0,
			Antenna:   0,
			Channel:   1,
			CodeRate:  "4/5",
			CrcStatus: 1,
			DataRate: DataRate{
				Bandwidth:    125,
				Modulation:   "LORA",
				SpreadFactor: 7,
			},
			Frequency: 868300000,
			LoRaSNR:   7,
			Mac:       g.Gateway.Mac,
			RfChain:   1,
			Rssi:      -57,
			Size:      23,
		},
		PhyPayload: phyBytes,
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
