package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/roadmap", func(c *gin.Context) {
		c.HTML(200, "roadmap.html", nil)
	})

	router.GET("/appearance", func(c *gin.Context) {
		c.HTML(200, "appearance.html", nil)
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
