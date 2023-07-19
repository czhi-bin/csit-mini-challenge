package main

import (
	routes "csit-mini-challenge/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/flight", routes.GetCheapestFlights)
	router.GET("/hotel", routes.GetCheapestHotels)

	router.Run(":8080")
}