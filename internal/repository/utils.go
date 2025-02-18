package repository

import (
	"context"
	"time"
)

// getContextWithTimeout helps as utility to pass on the timeout controlled context.
func getContextWithTimeout(timeoutInSecond int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeoutInSecond)*time.Second)
}
