package accrual

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_Fire(t *testing.T) {
	r := NewRateLimiter(60, time.Second)
	for i := 0; i < 60; i++ {
		assert.Less(t, 0, r.Fire())
	}
	require.Equal(t, 0, r.Fire())
	time.Sleep(time.Second)
	require.Equal(t, 60, r.Fire())
}
