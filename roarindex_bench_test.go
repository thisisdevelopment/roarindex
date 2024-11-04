package roarindex

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkRoarIndexPushMap(b *testing.B) {
	om := NewRoarIndex[string, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapID := fmt.Sprintf("map%d", i%1000)
		value := i % 100
		om.PushMap(mapID, value)
	}
}

func BenchmarkRoarIndexGetMap(b *testing.B) {
	om := NewRoarIndex[string, int]()
	for i := 0; i < 1000; i++ {
		mapID := fmt.Sprintf("map%d", i)
		for j := 0; j < 10; j++ {
			om.PushMap(mapID, j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapID := fmt.Sprintf("map%d", i%1000)
		om.GetMap(mapID)
	}
}

func BenchmarkRoarIndexHasValue(b *testing.B) {
	om := NewRoarIndex[string, int]()
	for i := 0; i < 1000; i++ {
		mapID := fmt.Sprintf("map%d", i)
		for j := 0; j < 10; j++ {
			om.PushMap(mapID, j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapID := fmt.Sprintf("map%d", i%1000)
		value := i % 10
		om.HasValue(mapID, value)
	}
}

func BenchmarkRoarIndexDeleteMap(b *testing.B) {
	om := NewRoarIndex[string, int]()
	for i := 0; i < 1000; i++ {
		mapID := fmt.Sprintf("map%d", i)
		for j := 0; j < 10; j++ {
			om.PushMap(mapID, j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapID := fmt.Sprintf("map%d", i%1000)
		om.DeleteMap(mapID)
	}
}

func BenchmarkMemoryUsageComparison(b *testing.B) {
	const numKeys = 10_000
	const numValues = 50_000

	b.Run(fmt.Sprintf("RoarIndex-%d-kv", numKeys*numValues), func(b *testing.B) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		beforeAlloc := m.Alloc

		om := NewRoarIndex[string, int]()
		// Insert 10,000 keys with ~50,000 values each (500,000,000 total values)
		for i := 0; i < numKeys; i++ {
			key := fmt.Sprintf("key%d", i)
			for j := 0; j < numValues; j++ {
				om.PushMap(key, j)
			}
		}

		runtime.ReadMemStats(&m)
		b.ReportMetric(float64(m.Alloc-beforeAlloc)/1024/1024, "MB-RoarIndex")
	})

	b.Run(fmt.Sprintf("StandardMap-%d-kv", numKeys*numValues), func(b *testing.B) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		beforeAlloc := m.Alloc

		standardMap := make(map[string][]int)
		// Insert 10,000 keys with ~50,000 values each (500,000,000 total values)
		for i := 0; i < numKeys; i++ {
			key := fmt.Sprintf("key%d", i)
			values := make([]int, 0, 5)
			for j := 0; j < numValues; j++ {
				values = append(values, j)
			}
			standardMap[key] = values
		}

		runtime.ReadMemStats(&m)
		b.ReportMetric(float64(m.Alloc-beforeAlloc)/1024/1024, "MB-StdMap")
	})
}

// ... existing code ...
