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

package generator

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/toskatok/lg/generator/atrovan"
	"github.com/toskatok/lg/generator/fanco"
	"github.com/toskatok/lg/generator/isrc"
	"github.com/toskatok/lg/generator/json"
	"github.com/toskatok/lg/generator/lora"
	"github.com/toskatok/lg/generator/ttn"
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
	case "isrc": // generators/isrc/isrc.go
		var isrc isrc.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &isrc); err != nil {
			return nil, err
		}

		return isrc, nil
	case "atrovan": // generators/atrovan/atrovan.go
		var atrovan atrovan.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &atrovan); err != nil {
			return nil, err
		}

		return atrovan, nil
	case "fanco": // generators/fanco/fanco.go
		var fanco fanco.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &fanco); err != nil {
			return nil, err
		}

		return fanco, nil
	case "ttn": // generators/ttn/ttn.go
		var ttn ttn.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &ttn); err != nil {
			return nil, err
		}

		return ttn, nil
	case "lora": // generators/lora/lora.go
		var lora lora.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &lora); err != nil {
			return nil, err
		}

		return lora, nil
	case "json": // generators/json/json.go
		var json json.Generator
		// load genrator information from configuration file
		if err := mapstructure.Decode(config, &json); err != nil {
			return nil, err
		}

		return json, nil
	}

	return nil, fmt.Errorf("generator %s is not supported yet", name)
}
