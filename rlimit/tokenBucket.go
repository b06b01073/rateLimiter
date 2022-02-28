package rlimit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type token struct{}

type bucket struct {
	bucketSize int
	fillRate   int
	tokens     chan token
}

var BucketSize int = 4
var FillRate int = 2
var B *bucket

func init() {
	B = initializeBucket(BucketSize, FillRate)
	go B.ticking()
}

func initializeBucket(bucketSize int, fillRate int) *bucket {
	b := &bucket{
		bucketSize: bucketSize,
		fillRate:   fillRate,
		tokens:     make(chan token, bucketSize),
	}

	return b
}

func (b *bucket) ticking() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			b.fillBucket()
		}
	}
}

func (b *bucket) fillBucket() {
	for i := 0; i < b.fillRate; i++ {
		select {
		case b.tokens <- token{}:
		default:
		}
	}
}

func (b *bucket) consumeToken() bool {
	select {
	case <-b.tokens:
		return true
	default:
	}
	return false
}

func TokenBucket(c *gin.Context) {
	if accepted := B.consumeToken(); !accepted {
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
	c.Next()
}
