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
	"fmt"
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

// List returns all running instances This function is mapped
// to the path GET /instances
func (v InstancesResource) List(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(instances))
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
	go func() {
		socket.BroadcastTo("I1820", req.Name, i.R.Count())
	}()
	instances[req.Name] = i

	return c.Render(http.StatusOK, r.JSON(true))
}

// Show shows the detail of given instance. This function is mapped
// to the path GET /instances/{instance_id}
func (v InstancesResource) Show(c buffalo.Context) error {
	id := c.Param("instance_id")
	i, ok := instances[id]
	if !ok {
		return c.Error(http.StatusNotFound, fmt.Errorf("There is no instnace with name %s", id))
	}

	return c.Render(http.StatusOK, r.JSON(i.R.Count()))
}

// Destroy stops given instance and removes it from the instances list.
// This function is mapped to the path DELETE /instances/{instance_id}
func (v InstancesResource) Destroy(c buffalo.Context) error {
	id := c.Param("instance_id")
	i, ok := instances[id]
	if !ok {
		return c.Error(http.StatusNotFound, fmt.Errorf("There is no instnace with name %s", id))
	}

	i.Stop()
	delete(instances, id)

	return c.Render(http.StatusOK, r.JSON(true))
}
