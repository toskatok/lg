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

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/toskatok/lg/models"
)

func (suite *LGTestSuite) Test_InstancesResource_Create() {
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

	w := httptest.NewRecorder()
	data, err := json.Marshal(instanceReq{
		Name:   "The whom was not given",
		Config: config,
	})
	suite.NoError(err)
	req, err := http.NewRequest("POST", "/api/instances", bytes.NewReader(data))
	suite.NoError(err)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
	suite.Contains(w.Body.String(), true)
}
