package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "message": "hello world",
  })
}

func main() {
  router := gin.Default()
  router.GET("/hello", func(c *gin.Context) {
    c.String(http.StatusOK, "hello world")
  })
  router.GET("/", helloHandler)

  router.Run("localhost:8080")
}