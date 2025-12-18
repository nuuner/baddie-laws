package main

import "github.com/gin-gonic/gin"

func main() {
  router := gin.Default()
  router.LoadHTMLGlob("templates/*")

  router.GET("/", func(c *gin.Context) {
    c.HTML(200, "index.html", nil)
  })

  router.Run() // listens on 0.0.0.0:8080 by default
}
