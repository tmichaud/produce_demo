package router

import (
	"produce_demo/api"

	"github.com/labstack/echo"
)

// Create a new Echo and add the api routes
func New() *echo.Echo {
	e := echo.New()

	// set main routes
	api.Produce(e)

	return e
}
