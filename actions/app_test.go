package actions

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// LGTestSuite is a test suite for lg APIs.
type LGTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates lg test suite
func (suite *LGTestSuite) SetupSuite() {
	suite.engine = App()
}

// Let's test lg APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(LGTestSuite))
}
