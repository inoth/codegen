package test

import (
	"testing"
	"time"
)

func TestQuickStar(t *testing.T) {
	t1 := time.Now()
	t.Logf("ok; time: %v", time.Since(t1))
}
