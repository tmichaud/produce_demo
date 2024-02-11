package router

import (
	"example.com/produce_demo/api"

	"github.com/labstack/echo/v4"
)

// Create a new Echo and add the api routes
func New() *echo.Echo {
	e := echo.New()

	// set main routes
	api.Produce(e)

	return e
}
