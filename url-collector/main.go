package main

import (
	"os"
	"url-collector/url-collector/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	defaultPort = "8080"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	picturesController := api.NewPicturesController()

	router.GET("/pictures", picturesController.GetImages)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router.Run(":" + port)
}
