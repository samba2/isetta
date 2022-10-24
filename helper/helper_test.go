package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryExhausted(t *testing.T) {
	err := Retry(RetryParams{
		Description: "foo",
		Attempts: 1,
		Sleep: 1 * time.Millisecond,
		Func: func() bool {return false},
	})

	assert.Error(t, err, "error of exhausted retries")
}

func TestRetrySuccessful(t *testing.T) {
	cnt := 0

	retryableFunc := func() bool {
		if cnt == 0 {
			cnt++
			return false
		} else {
			return true
		}
	}

	err := Retry(RetryParams{
		Description: "retry was successful",
		Attempts: 2,
		Sleep: 1 * time.Millisecond,
		Func: retryableFunc,
	})

	assert.NoError(t, err, "retry was successful")
}