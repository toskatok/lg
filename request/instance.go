package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/toskatok/lg/instance"
)

// Instance creation request
type Instance struct {
	Name string
	instance.Config
}

func (r Instance) Validate() error {
	return validation.ValidateStruct(&r, validation.Field(&r.Name, validation.Required))
}
