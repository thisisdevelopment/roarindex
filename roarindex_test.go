package roarindex

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/thoas/go-funk"
)

func TestRoarIndexSetAndGet(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)

	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}

	// Check if the result contains the correct value
	if len(result) != 1 || result[0] != value {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetNonExisting(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)

	// Test GetMap (Get)
	result, err := om.GetMap(mapID + "nonExistent")
	// Check if the map exists
	if err == nil {
		t.Errorf("Expected map to not exist, but it does")
	}

	// Check if the result contains the correct value
	if len(result) != 0 {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}
}

func TestRoarIndexHasValueNonExisting(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)

	// Test GetMap (Get)
	result := om.HasValue(mapID+"nonExistent", 42)

	// Check if the result contains the correct value
	if result {
		t.Errorf("Expected result to be false, but got true")
	}
}

func TestRoarIndexSetAndGetDups(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)
	om.PushMap(mapID, value)
	om.PushMap(mapID, value)
	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}
	// Check if the result contains the correct value
	if len(result) != 1 || result[0] != value {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}
	if len(result) != len(om.Values()) {
		t.Errorf("Expected result to be %d, but got %d", len(om.Values()), len(result))
	}

	if len(om.Keys()) != 1 {
		t.Errorf("Expected result to be %d, but got %d", 1, len(om.Keys()))
	}
	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetDups2(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)
	om.PushMap(mapID, value)
	om.PushMap(mapID, value)
	om.PushMap(mapID, value+1)
	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}
	// Check if the result contains the correct value
	if len(result) != 2 || result[0] != value {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}
	if len(result) != len(om.Values()) {
		t.Errorf("Expected result to be %d, but got %d", len(om.Values()), len(result))
	}

	if len(om.Keys()) != 1 {
		t.Errorf("Expected result to be %d, but got %d", 1, len(om.Keys()))
	}
	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetMultiple(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	value2 := 43
	om.PushMap(mapID, value)

	om.PushMap(mapID, value2)

	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}

	// Check if the result contains the correct value
	if len(result) != 2 || result[0] != value {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}
	if len(result) != 2 || result[1] != value2 {
		t.Errorf("Expected result to be [%d], but got %v", value2, result)
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetMultipleDifferentMaps(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	mapID2 := "testMap2"
	value := 42
	value2 := 43
	om.PushMap(mapID, value)

	om.PushMap(mapID2, value2)

	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}

	// Check if the result contains the correct value
	if len(result) != 1 || result[0] != value {
		t.Errorf("Expected result to be [%d], but got %v", value, result)
	}
	result2, err := om.GetMap(mapID2)
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't")
	}
	if len(result2) != 1 || result2[0] != value2 {
		t.Errorf("Expected result to be [%d], but got %v", value2, result2)
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetMultipleSingleMapGrow(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Test PushMap (Set)
	mapID := "testMap"
	value := 42
	om.PushMap(mapID, value)
	om.PushMap(mapID, value+1)
	om.PushMap(mapID, value+2)
	om.PushMap(mapID, value+3)

	// Test GetMap (Get)
	result, err := om.GetMap(mapID)
	// Check if the map exists
	if err != nil {
		t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
	}

	// Check if the result contains the correct value
	if len(result) != 4 || result[0] != value || result[1] != value+1 || result[2] != value+2 || result[3] != value+3 {
		t.Errorf("Expected result to be [%d, %d, %d, %d], but got %v", value, value+1, value+2, value+3, result)
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexSetAndGetMultipleMultiMapGrow(t *testing.T) {
	// Initialize the RoarIndex
	maps := []string{"testMap1", "testMap2", "testMap3"}
	om := NewRoarIndex[string, int]()

	for _, mapID := range maps {

		// Test PushMap (Set)
		value := 42
		om.PushMap(mapID, value)
		om.PushMap(mapID, value+1)
		om.PushMap(mapID, value+2)
		om.PushMap(mapID, value+3)

		// Test GetMap (Get)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
		}

		// Check if the result contains the correct value
		if len(result) != 4 || result[0] != value || result[1] != value+1 || result[2] != value+2 || result[3] != value+3 {
			t.Errorf("Expected result to be [%d, %d, %d, %d], but got %v", value, value+1, value+2, value+3, result)
		}

		// Test GetMap for a non-existent map
		nonExistentResult, err := om.GetMap("nonExistentMap")

		// Check if the non-existent map is correctly reported as not found
		if err == nil {
			t.Errorf("Expected non-existent map to return false, but got true")
		}

		// Check if the result for non-existent map is an empty slice
		if len(nonExistentResult) != 0 {
			t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
		}
	}
}

func TestRoarIndexSetAndGetMultipleMultiMapGrowSameValues(t *testing.T) {
	// Initialize the RoarIndex
	maps := []string{"testMap1", "testMap2", "testMap3"}
	om := NewRoarIndex[string, int]()

	for _, mapID := range maps {

		// Test PushMap (Set)
		value := 42
		om.PushMap(mapID, value)
		om.PushMap(mapID, value+1)
		om.PushMap(mapID, value+2)
		om.PushMap(mapID, value+3)

		// Test GetMap (Get)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
		}

		// Check if the result contains the correct value
		if len(result) != 4 || result[0] != 42 || result[1] != 43 || result[2] != 44 || result[3] != 45 {
			t.Errorf("Expected result to be [%d, %d, %d, %d], but got %v", 42, 43, 44, 45, result)
		}

		// Test GetMap for a non-existent map
		nonExistentResult, err := om.GetMap("nonExistentMap")

		// Check if the non-existent map is correctly reported as not found
		if err == nil {
			t.Errorf("Expected non-existent map to return false, but got true")
		}

		// Check if the result for non-existent map is an empty slice
		if len(nonExistentResult) != 0 {
			t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
		}
	}
}

func TestRoarIndexSetAndGetMultipleMultiMapGrowDiffValues(t *testing.T) {
	// Initialize the RoarIndex
	maps := []string{"testMap1", "testMap2", "testMap3"}
	om := NewRoarIndex[string, int]()

	for i, mapID := range maps {

		// Test PushMap (Set)
		value := i * 10
		om.PushMap(mapID, value)
		om.PushMap(mapID, value+1)
		om.PushMap(mapID, value+2)
		om.PushMap(mapID, value+3)

		// Test GetMap (Get)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map to exist, but it doesn't, error: %s", err.Error())
		}
		slices.Sort(result)
		// Check if the result contains the correct value
		if len(result) != 4 || result[0] != i*10 || result[1] != i*10+1 || result[2] != i*10+2 || result[3] != i*10+3 {
			t.Errorf("Expected result to be [%d, %d, %d, %d], but got %v", i*10, i*10+1, i*10+2, i*10+3, result)
		}

		// Test GetMap for a non-existent map
		nonExistentResult, err := om.GetMap("nonExistentMap")

		// Check if the non-existent map is correctly reported as not found
		if err == nil {
			t.Errorf("Expected non-existent map to return false, but got true")
		}

		// Check if the result for non-existent map is an empty slice
		if len(nonExistentResult) != 0 {
			t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
		}
	}
}

// ... existing code ...

func TestRoarIndexLargeSetAndGet(t *testing.T) {
	const (
		amountRuns   = 1
		amountKeys   = 100000
		amountValues = 500
	)

	// Initialize the RoarIndex with a larger capacity
	om := NewRoarIndex[string, int]()
	controlPool := make(map[string][]int)

	// Generate 1000 keys and values
	for y := 0; y < amountRuns; y++ {
		t.Logf("run %d", y)
		for i := 0; i < amountKeys; i++ {
			mapID := fmt.Sprintf("map%d", i)
			value := rand.Intn(amountValues) // Random value between 0 and 49 (allowing duplicates)
			om.PushMap(mapID, value)
			controlPool[mapID] = funk.UniqInt(append(controlPool[mapID], value))
		}
	}

	// Test GetMap for all 100 keys
	for i := 0; i < amountKeys; i++ {
		mapID := fmt.Sprintf("map%d", i)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map %s to exist, but it doesn't, error: %s", mapID, err.Error())
		}

		// Check if the result contains at least one value
		if len(result) == 0 {
			t.Errorf("Expected result for map %s to contain at least one value, but it's empty", mapID)
		}

		slices.Sort(controlPool[mapID])
		slices.Sort(result)
		if !reflect.DeepEqual(controlPool[mapID], result) {
			t.Errorf("Expected value in map %s to be %v, but got %v", mapID, controlPool[mapID], result)
		}

		// Check if all values are within the expected range
		for _, v := range result {
			if v < 0 || v >= amountValues {
				t.Errorf("Expected value in map %s to be between 0 and %d, but got %d", mapID, amountValues-1, v)
			}

			if !om.HasValue(mapID, v) {
				t.Errorf("Expected value %d to be in map %s", v, mapID)
			}
		}
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexLargeSetAndGetStrings(t *testing.T) {
	const (
		amountRuns   = 1
		amountKeys   = 1000000
		amountValues = 500000
	)

	// Initialize the RoarIndex with a larger capacity
	om := NewRoarIndex[string, string]()

	controlPool := make(map[string][]string)

	// Generate 1000 keys and values
	for y := 0; y < amountRuns; y++ {
		t.Logf("run %d", y)
		for i := 0; i < amountKeys; i++ {
			mapID := fmt.Sprintf("map%d", i)
			// Generate a random string of 128 bytes
			randomBytes := make([]byte, 128)
			_, err := crand.Read(randomBytes)
			if err != nil {
				t.Fatalf("Failed to generate random bytes: %v", err)
			}
			randomString := string(randomBytes)

			value := fmt.Sprintf("%s%d", randomString, rand.Intn(amountValues))
			om.PushMap(mapID, value)
			controlPool[mapID] = funk.UniqString(append(controlPool[mapID], value))
		}
	}

	// Test GetMap for all 100 keys
	for i := 0; i < amountKeys; i++ {
		mapID := fmt.Sprintf("map%d", i)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map %s to exist, but it doesn't, error: %s", mapID, err.Error())
		}

		// Check if the result contains at least one value
		if len(result) == 0 {
			t.Errorf("Expected result for map %s to contain at least one value, but it's empty", mapID)
		}

		slices.Sort(controlPool[mapID])
		slices.Sort(result)
		if !reflect.DeepEqual(controlPool[mapID], result) {
			t.Errorf("Expected value in map %s to be %v, but got %v", mapID, controlPool[mapID], result)
		}

		// Check if all values are within the expected range
		for _, v := range result {
			if !om.HasValue(mapID, v) {
				t.Errorf("Expected value %s to be in map %s", v, mapID)
			}
		}
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexLargeSetAndGetStructs(t *testing.T) {
	const (
		amountRuns   = 10
		amountKeys   = 100000
		amountValues = 5000
	)

	type TestStruct struct {
		Name        string
		RandomValue string
	}

	// Initialize the RoarIndex with a larger capacity
	om := NewRoarIndex[string, TestStruct]()

	controlPool := make(map[string][]TestStruct)

	// Generate 1000 keys and values
	for y := 0; y < amountRuns; y++ {
		t.Logf("run %d", y)
		for i := 0; i < amountKeys; i++ {
			mapID := fmt.Sprintf("map%d", i)
			// Generate a random string of 128 bytes
			randomBytes := make([]byte, 128)
			_, err := crand.Read(randomBytes)
			if err != nil {
				t.Fatalf("Failed to generate random bytes: %v", err)
			}
			randomString := string(randomBytes)

			value := TestStruct{
				Name:        randomString,
				RandomValue: fmt.Sprintf("%s%d", randomString, rand.Intn(amountValues)),
			}
			om.PushMap(mapID, value)
			controlPool[mapID] = append(controlPool[mapID], value)
		}
	}

	// Test GetMap for all 100 keys
	for i := 0; i < amountKeys; i++ {
		mapID := fmt.Sprintf("map%d", i)
		result, err := om.GetMap(mapID)
		// Check if the map exists
		if err != nil {
			t.Errorf("Expected map %s to exist, but it doesn't, error: %s", mapID, err.Error())
		}

		// Check if the result contains at least one value
		if len(result) == 0 {
			t.Errorf("Expected result for map %s to contain at least one value, but it's empty", mapID)
		}

		if !reflect.DeepEqual(controlPool[mapID], result) {
			t.Errorf("Expected value in map %s to be %v, but got %v", mapID, controlPool[mapID], result)
		}

		// Check if all values are within the expected range
		for _, v := range result {
			if !om.HasValue(mapID, v) {
				t.Errorf("Expected value %s to be in map %s", v, mapID)
			}
		}
	}

	// Test GetMap for a non-existent map
	nonExistentResult, err := om.GetMap("nonExistentMap")

	// Check if the non-existent map is correctly reported as not found
	if err == nil {
		t.Errorf("Expected non-existent map to return false, but got true")
	}

	// Check if the result for non-existent map is an empty slice
	if len(nonExistentResult) != 0 {
		t.Errorf("Expected empty slice for non-existent map, but got %v", nonExistentResult)
	}
}

func TestRoarIndexDelete(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Add some test data
	om.PushMap("map1", 1)
	om.PushMap("map2", 2)
	om.PushMap("map3", 3)
	om.PushMap("map4", 4)

	// Test deleting an existing map
	om.DeleteMap("map1")

	// Check if the deleted map no longer exists
	result, err := om.GetMap("map1")
	if err == nil {
		t.Errorf("Expected map1 to be deleted, but it still exists")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result for deleted map, but got %v", result)
	}

	// Check if other maps are unaffected
	result2, err := om.GetMap("map2")
	if err != nil {
		t.Errorf("Expected map2 to exist, but it doesn't, error: %s", err.Error())
	}
	if !reflect.DeepEqual(result2, []int{2}) {
		t.Errorf("Expected map2 to contain [2], but got %v", result2)
	}

	// Test deleting a non-existent map (should not cause any errors)
	om.DeleteMap("nonExistentMap")

	// Verify that existing maps are still intact
	result3, err := om.GetMap("map3")
	if err != nil {
		t.Errorf("Expected map3 to exist, but it doesn't, error: %s", err.Error())
	}
	if !reflect.DeepEqual(result3, []int{3}) {
		t.Errorf("Expected map3 to contain [3], but got %v", result3)
	}
}

func TestRoarIndexMultiPush(t *testing.T) {
	// Initialize the RoarIndex
	om := NewRoarIndex[string, int]()

	// Add some test data
	om.PushMap("map1", 1)
	om.PushMap("map1", 1)
	om.PushMap("map1", 1)
	om.PushMap("map1", 1)
	om.PushMap("map1", 1)

	// Test deleting an existing map
	lengthValues := len(om.Values())
	if lengthValues != 1 {
		t.Errorf("Expected map1 to have only one value, but it has %d", lengthValues)
	}

	lengthKeys := om.Count()
	if lengthKeys != 1 {
		t.Errorf("Expected map1 to have only one key, but it has %d", lengthKeys)
	}
}

func TestRoarIndexConcurrentAccess(t *testing.T) {
	om := NewRoarIndex[string, string]()
	var wg sync.WaitGroup

	keys := []string{"key1", "key2", "key3"}
	values := []string{"value1", "value2", "value3"}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			value := values[i%len(values)]
			om.PushMap(key, value)
		}(i)
	}

	wg.Wait()

	for _, key := range keys {
		_, err := om.GetMap(key)
		if err != nil && err != ErrKeyNotFound {
			t.Errorf("GetMap failed for key %q: %v", key, err)
		}
		// Optionally, verify the values if needed
	}
}

func TestRoarIndexRandomized(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	om := NewRoarIndex[int, int]()
	numOps := 1000

	for i := 0; i < numOps; i++ {
		opType := rand.Intn(3)
		key := rand.Intn(100)
		value := rand.Intn(1000)
		switch opType {
		case 0:
			om.PushMap(key, value)
		case 1:
			_, _ = om.GetMap(key)
		case 2:
			om.DeleteMap(key)
		}
	}
}
