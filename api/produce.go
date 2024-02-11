package api

import (
	"example.com/produce_demo/api/handlers"

	"github.com/labstack/echo/v4"
)

func Produce(e *echo.Echo) {
	// Add a new Produce item to Inventory
	e.POST("/produce", handlers.AddProduce)

	// Delete Produce item from Inventory
	e.DELETE("/produce/:ProduceCode", handlers.DeleteProduce)

	// Fetch all Produce items from Inventory
	e.GET("/produce", handlers.FetchProduce)

	// Fetch a Produce item from Inventory by Produce Code
	e.GET("/produce/:ProduceCode", handlers.FetchProduceByProduceCode)
}
