/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-07-2018
 * |
 * | File Name:     generator.go
 * +===============================================
 */

package generators

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Generator generates data whenever you want
// based on given input. input can be nil.
type Generator interface {
	Generate(input interface{}) ([]byte, error)
	Topic() string
}

// Get return a configured instance of Generator based on given name
func Get(name string, config interface{}) (Generator, error) {
	// Generator selection and configuration
	switch name {
	case "isrc": // generators/isrc.go
		var isrc ISRCGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &isrc); err != nil {
			return nil, err
		}

		return isrc, nil
	case "atrovan": // generators/atrovan.go
		var atrovan AtrovanGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &atrovan); err != nil {
			return nil, err
		}

		return atrovan, nil
	case "fanco": // generators/fanco.go
		var fanco FancoGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &fanco); err != nil {
			return nil, err
		}

		return fanco, nil
	case "ttn": // generators/ttn.go
		var ttn TTNGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &ttn); err != nil {
			return nil, err
		}

		return ttn, nil
	case "lora": // generators/lora.go
		var lora LoRaGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &lora); err != nil {
			return nil, err
		}

		return lora, nil
	case "json": // generators/json.go
		var json JSONGenerator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &json); err != nil {
			return nil, err
		}

		return json, nil
	}

	return nil, fmt.Errorf("generator %s is not supported yet", name)
}
