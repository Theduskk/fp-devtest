package main

import (
	apis "flatpeak-devtask/apis"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/slots", apis.GetSlots)
	router.Run("0.0.0.0:3000")
}
