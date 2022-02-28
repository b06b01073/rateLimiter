package rlimit

import (
	"net/http"
	"testing"
	"time"
)

var testing_URL = "http://localhost:8080/ping"

func TestTokenBucket(t *testing.T) {
	tests := []struct {
		requests int
	}{
		{1},
		{2},
		{BucketSize},
		{BucketSize + 3},
	}

	for _, test := range tests {
		start := time.Now()
		requests := test.requests

		var flag bool
		for i := 1; i < requests; i++ {
			resp, err := http.Get(testing_URL)
			if err != nil {
				t.Errorf(err.Error())
			}
			if requests <= BucketSize && resp.StatusCode == http.StatusTooManyRequests {
				t.Errorf("Reject request when bucket is not empty")
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				flag = true
			}
		}

		if !flag && requests > BucketSize {
			elapse := time.Since(start)
			t.Errorf("Accept request when bucket is empty. %f", elapse.Seconds())
		}
	}
}
