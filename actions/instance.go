/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-12-2018
 * |
 * | File Name:     instance.go
 * +===============================================
 */

package actions

import (
	"net/http"
	"time"

	"github.com/I1820/lg/models"
	"github.com/gobuffalo/buffalo"
)

// instance creation request
type instanceReq struct {
	Name string
	models.Config
}

// list of the running instances
var instances map[string]*models.Instance = make(map[string]*models.Instance)

// InstancesResource manages instances of load generators.
type InstancesResource struct {
	buffalo.Resource
}

// List returns all running instances
func (v InstancesResource) List(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON([]int{}))
}

// Create runs new generator instance. This function is mapped
// to the path POST /instances
func (v InstancesResource) Create(c buffalo.Context) error {
	var req instanceReq
	if err := c.Bind(&req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	var config models.Config = req.Config

	rate, err := time.ParseDuration(c.Param("rate"))
	if err != nil {
		rate = 1 * time.Millisecond
	}
	destination := c.Param("destination")
	if destination == "" {
		destination = "mqtt://127.0.0.1:1883"
	}

	i, err := models.NewInstance(config, rate, destination)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	i.Run()
	instances[req.Name] = i

	return c.Render(http.StatusOK, r.JSON(true))
}
