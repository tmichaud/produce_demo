package main

import (
	"fmt"
	router "produce_demo/routers"
)

// Main Function
func main() {
	fmt.Println("Welcome to the webserver")
	e := router.New()
	e.Start(":8080")
}
