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

	"github.com/labstack/echo/v4"
	"github.com/toskatok/lg/models"
)

func (suite *LGTestSuite) Test_InstancesResource_Create() {
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

	cw := httptest.NewRecorder()
	data, err := json.Marshal(instanceReq{
		Name:   "kj",
		Config: config,
	})
	suite.NoError(err)
	creq, err := http.NewRequest("POST", "/api/instances", bytes.NewReader(data))
	suite.NoError(err)
	creq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	q := creq.URL.Query()
	q.Add("destination", "http://127.0.0.1:8080")
	q.Add("rate", "1m")
	creq.URL.RawQuery = q.Encode()
	suite.engine.ServeHTTP(cw, creq)

	suite.Equal(200, cw.Code)
	suite.NoError(json.Unmarshal(cw.Body.Bytes(), &t))
	suite.Equal(t, true)

	dw := httptest.NewRecorder()
	dreq, err := http.NewRequest("DELETE", "/api/instances/kj", nil)
	suite.NoError(err)
	suite.engine.ServeHTTP(dw, dreq)

	suite.Equal(200, dw.Code)
	suite.NoError(json.Unmarshal(dw.Body.Bytes(), &t))
	suite.Equal(t, true)
}
