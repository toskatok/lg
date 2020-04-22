package request

import "github.com/toskatok/lg/instance"

// Instance creation request
type Instance struct {
	Name string
	instance.Config
}
