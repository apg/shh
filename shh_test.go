package main

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkSendingMeasurements(b *testing.B) {
	measurements := make(chan Measurement, 100)
	tick := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for _ = range measurements {
		}
		wg.Done()
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		measurements <- Measurement{tick, "test", []string{"testing", "thing"}, 1.0}
	}
	close(measurements)
	wg.Wait()
}
