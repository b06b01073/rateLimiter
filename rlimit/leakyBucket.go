package rlimit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type request chan interface{}
type signal struct{}

type queue struct {
	queueSize    int
	outflowRate  int
	requestQueue chan string
	signals      map[string]*request
}

var queueSize int = 4
var outflowRate int = 2
var Q *queue

func init() {
	Q = initializeQueue()
	go Q.ticking()
}

func initializeQueue() *queue {
	return &queue{
		queueSize:    queueSize,
		outflowRate:  outflowRate,
		requestQueue: make(chan string, queueSize),
		signals:      make(map[string]*request),
	}
}

func (q *queue) ticking() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		q.outflow()
	}
}

func (q *queue) outflow() {
	var u string
	s := signal{}
	for i := 0; i < q.outflowRate; i++ {
		select {
		case u = <-q.requestQueue:
			*(q.signals[u]) <- s
			delete(q.signals, u)
		default:
		}
	}
}

func (q *queue) enqueue(r request, uuid string) bool {
	select {
	case Q.requestQueue <- uuid:
		q.signals[uuid] = &r
		return true
	default:
		return false
	}
}

func LeakyBucket(c *gin.Context) {
	/*
		if queue is full:
			reject the request
		else
			request enters the queue and blocked until it is fetched
	*/

	// generate a uuid for request
	uuid := uuid.NewString()
	r := make(chan interface{})
	ok := Q.enqueue(r, uuid)

	if !ok {
		c.AbortWithStatus(http.StatusTooManyRequests)
	}

	// gin的abort會執行完當前handler才停止，如果狀態為ok才會需要block住
	if ok {
		<-r
	}

	c.Next()
}
