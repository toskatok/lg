/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-12-2018
 * |
 * | File Name:     instance_test.go
 * +===============================================
 */

package actions

import "github.com/toskatok/lg/models"

func (as *ActionSuite) Test_InstancesResource_Create() {
	var t bool

	// KJ Configurations
	var config models.Config
	config.Generator.Name = "ttn"
	config.Token = "ttnIStheBEST"
	config.Messages = []map[string]interface{}{
		{
			"101": 6000,
			"100": 6500,
		},
	}

	// creates new instance
	res := as.JSON("/api/instances").Post(instanceReq{
		Name:   "The whom was not given",
		Config: config,
	})
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&t)
	as.Equal(t, true)
}
