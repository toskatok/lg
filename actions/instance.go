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

	"github.com/gobuffalo/buffalo"
)

// InstancesResource manages instances of load generators.
type InstancesResource struct {
	buffalo.Resource
}

// List returns all running instances
func (v InstancesResource) List(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON([]int{}))
}
