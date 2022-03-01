package main

import (
	"rateLimiter/rlimit"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/tokenBucket", rlimit.TokenBucket, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "request accepted",
		})
	})

	r.GET("/leakyBucket", rlimit.LeakyBucket, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "request accepted",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
