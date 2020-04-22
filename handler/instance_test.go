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

package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/toskatok/lg/handler"
	"github.com/toskatok/lg/instance"
	"github.com/toskatok/lg/request"
)

// LGTestSuite is a test suite for lg APIs.
type LGTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates lg test suite
func (suite *LGTestSuite) SetupSuite() {
	suite.engine = echo.New()
	handler.NewInstance().Register(suite.engine.Group("/api"))
}

// Let's test lg APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(LGTestSuite))
}

const (
	value100 = 6000
	value101 = 6500
	name     = "elahe"
)

func (suite *LGTestSuite) Test_InstancesResource_Create() {
	var t bool

	// Elahe instance configurations
	var config instance.Config
	config.Generator.Name = "ttn"
	config.Token = "ttnIStheBEST"
	config.Messages = []map[string]interface{}{
		{
			"101": value100,
			"100": value101,
		},
	}

	cw := httptest.NewRecorder()
	data, err := json.Marshal(request.Instance{
		Name:   name,
		Config: config,
	})
	suite.NoError(err)

	creq := httptest.NewRequest("POST", "/api/instances", bytes.NewReader(data))
	creq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	q := creq.URL.Query()
	q.Add("destination", "http://127.0.0.1:8080")
	q.Add("rate", "1m")
	creq.URL.RawQuery = q.Encode()
	suite.engine.ServeHTTP(cw, creq)

	suite.Equal(http.StatusOK, cw.Code)
	suite.NoError(json.Unmarshal(cw.Body.Bytes(), &t))
	suite.Equal(t, true)

	dw := httptest.NewRecorder()
	dreq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/instances/%s", name), nil)
	suite.engine.ServeHTTP(dw, dreq)

	suite.Equal(http.StatusOK, dw.Code)
	suite.NoError(json.Unmarshal(dw.Body.Bytes(), &t))
	suite.Equal(t, true)
}
