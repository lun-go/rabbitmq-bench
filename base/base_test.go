package base

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkHello(b *testing.B) {
	b.ResetTimer()
	body := "1234567890"
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%s||%v", body, time.Now().UnixNano())
	}
}
