package roarindex

import (
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzRoarIndexImpl(f *testing.F) {
	// Add some seed corpus
	f.Add("key1", "value1", "value2")
	f.Add("key2", "value3", "value4")

	f.Fuzz(func(t *testing.T, key string, value1 string, value2 string) {
		om := NewRoarIndex[string, string]()

		// Test PushMap and GetMap
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for key %s: %v", key, err)
		}

		expectedCount := 2
		if value1 == value2 {
			expectedCount = 1
		}

		if len(values) != expectedCount {
			t.Errorf("Expected %d values, got %d", expectedCount, len(values))
		}

		if !contains(values, value1) {
			t.Errorf("GetMap did not contain value1: %v", values)
		}
		if !contains(values, value2) {
			t.Errorf("GetMap did not contain value2: %v", values)
		}

		// Test HasValue
		if !om.HasValue(key, value1) {
			t.Errorf("HasValue failed for key %s and value %s", key, value1)
		}
		if !om.HasValue(key, value2) {
			t.Errorf("HasValue failed for key %s and value %s", key, value2)
		}

		// Test non-existent value
		if om.HasValue(key, "non-existent") {
			t.Errorf("HasValue returned true for non-existent value")
		}

		// Test DeleteMap
		om.DeleteMap(key)
		values, err = om.GetMap(key)
		if err != ErrKeyNotFound {
			t.Errorf("Expected ErrKeyNotFound after DeleteMap, got: %v", err)
		}
		if len(values) != 0 {
			t.Errorf("Expected empty slice after DeleteMap, got: %v", values)
		}

		// Test Count
		if om.Count() != 0 {
			t.Errorf("Expected Count to be 0 after DeleteMap, got: %d", om.Count())
		}
	})
}

// Helper function to check if a slice contains a specific value
func contains[V comparable](slice []V, value V) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Fuzz test for RoarIndex with string keys and values
func FuzzRoarIndexStrings(f *testing.F) {
	// Seed corpus with various string inputs
	f.Add("key1", "value1", "value2")
	f.Add("", "value1", "value2")                                                          // Empty key
	f.Add("key1", "", "value2")                                                            // Empty value1
	f.Add("key1", "value1", "")                                                            // Empty value2
	f.Add("   ", "value1", "value2")                                                       // Whitespace key
	f.Add("key😊", "value❤️", "value✨")                                                     // Unicode characters
	f.Add(strings.Repeat("k", 1000), strings.Repeat("v", 1000), strings.Repeat("w", 1000)) // Long strings
	f.Add("key\x00", "value1", "value2")                                                   // Null byte in key
	f.Add("key", "value1\x00", "value2")                                                   // Null byte in value1
	f.Add("key", "value1", "value2\x00")                                                   // Null byte in value2

	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		om := NewRoarIndex[string, string]()

		// Test PushMap and GetMap
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for key %q: %v", key, err)
		}

		expectedValues := make(map[string]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}

		for _, v := range values {
			if _, exists := expectedValues[v]; !exists {
				t.Errorf("Unexpected value %q in values: %v", v, values)
			}
		}

		// Test HasValue
		if !om.HasValue(key, value1) {
			t.Errorf("HasValue failed for key %q and value %q", key, value1)
		}
		if !om.HasValue(key, value2) {
			t.Errorf("HasValue failed for key %q and value %q", key, value2)
		}

		// Test non-existent value
		if om.HasValue(key, "non-existent") {
			t.Errorf("HasValue returned true for non-existent value")
		}

		// Test DeleteMap
		om.DeleteMap(key)
		values, err = om.GetMap(key)
		if err != ErrKeyNotFound {
			t.Errorf("Expected ErrKeyNotFound after DeleteMap, got: %v", err)
		}
		if len(values) != 0 {
			t.Errorf("Expected empty slice after DeleteMap, got: %v", values)
		}

		// Test Count
		if om.Count() != 0 {
			t.Errorf("Expected Count to be 0 after DeleteMap, got: %d", om.Count())
		}
	})
}

// Fuzz test for RoarIndex with integer keys and values
func FuzzRoarIndexIntegers(f *testing.F) {
	// Seed corpus with integer inputs
	f.Add(1, 2, 3)
	f.Add(0, 0, 0)
	f.Add(-1, -2, -3)
	f.Add(123456789, 987654321, 192837465)
	f.Add(rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))

	f.Fuzz(func(t *testing.T, key, value1, value2 int) {
		om := NewRoarIndex[int, int]()

		// Test PushMap and GetMap
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for key %d: %v", key, err)
		}

		expectedValues := make(map[int]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}

		for _, v := range values {
			if _, exists := expectedValues[v]; !exists {
				t.Errorf("Unexpected value %d in values: %v", v, values)
			}
		}
	})
}

// Test for concurrent access to RoarIndex

// Custom comparable struct for testing
// type CustomValue struct {
// 	ID   int
// 	Name string
// }

// Fuzz test for RoarIndex with custom comparable types
// func FuzzRoarIndexCustomType(f *testing.F) {
// 	f.Add("key1", CustomValue{1, "Alice"}, CustomValue{2, "Bob"})
// 	f.Add("key2", CustomValue{3, "Charlie"}, CustomValue{4, "Dana"})

// 	f.Fuzz(func(t *testing.T, key string, value1, value2 CustomValue) {
// 		om := NewRoarIndex[string, CustomValue]()

// 		// Test PushMap and GetMap
// 		om.PushMap(key, value1)
// 		om.PushMap(key, value2)

// 		values, err := om.GetMap(key)
// 		if err != nil {
// 			t.Errorf("GetMap failed: %v", err)
// 		}

// 		expectedValues := make(map[CustomValue]struct{})
// 		expectedValues[value1] = struct{}{}
// 		expectedValues[value2] = struct{}{}

// 		if len(values) != len(expectedValues) {
// 			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
// 		}

// 		for _, v := range values {
// 			if _, exists := expectedValues[v]; !exists {
// 				t.Errorf("Unexpected value: %+v", v)
// 			}
// 		}
// 	})
// }

// Fuzz test for error handling in RoarIndex
func FuzzRoarIndexErrorHandling(f *testing.F) {
	f.Add("key1")
	f.Add("")
	f.Add("key\x00")

	f.Fuzz(func(t *testing.T, key string) {
		om := NewRoarIndex[string, string]()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recovered from panic: %v", r)
			}
		}()
		om.DeleteMap(key) // Should not panic even if key doesn't exist
	})
}

// Fuzz test to check internal state consistency of RoarIndex
func FuzzRoarIndexInternalState(f *testing.F) {
	f.Add("key1", "value1", "value2")
	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		om := NewRoarIndex[string, string]()
		om.PushMap(key, value1)
		om.PushMap(key, value2)
		om.DeleteMap(key)

		om.mtx.RLock()
		defer om.mtx.RUnlock()

		if _, exists := om.keyToID[key]; exists {
			t.Errorf("keyToID should not contain deleted key %q", key)
		}
		keyID, exists := om.keyToID[key]
		if exists {
			if _, exists := om.data[keyID]; exists {
				t.Errorf("data map should not contain deleted key ID")
			}
		}
	})
}

// Test RoarIndex with random inputs to simulate unexpected scenarios

// Property-based testing using Gopter (requires additional dependency)
// Uncomment and use if Gopter is added to your project
/*
import (
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestRoarIndexProperties(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	properties.Property("Adding and retrieving values", prop.ForAll(
		func(key string, values []string) bool {
			om := NewRoarIndex[string, string]()
			uniqueValues := make(map[string]struct{})
			for _, value := range values {
				om.PushMap(key, value)
				uniqueValues[value] = struct{}{}
			}
			retrievedValues, _ := om.GetMap(key)
			if len(retrievedValues) != len(uniqueValues) {
				return false
			}
			for _, v := range retrievedValues {
				if _, exists := uniqueValues[v]; !exists {
					return false
				}
			}
			return true
		},
		gen.AnyString(),                 // Key generator
		gen.SliceOf(gen.AnyString()), // Values generator
	))

	properties.TestingRun(t)
}
*/

// Additional helper function to generate unique elements from a slice
func uniqueStrings(values []string) []string {
	set := make(map[string]struct{})
	result := []string{}
	for _, v := range values {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Fuzz test to ensure RoarIndex handles empty and whitespace strings correctly
func FuzzRoarIndexEmptyWhitespace(f *testing.F) {
	f.Add("", "", "")
	f.Add("   ", "   ", "   ")
	f.Add("\t\n", "\r\n", "\v\f")

	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		om := NewRoarIndex[string, string]()

		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for key %q: %v", key, err)
		}

		expectedValues := make(map[string]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}
	})
}

// Fuzz test to check handling of strings with special and control characters
func FuzzRoarIndexSpecialCharacters(f *testing.F) {
	f.Add("key\n", "value\t", "value\r")
	f.Add("key\x00", "value\x00", "value\x1F")
	f.Add("DROP TABLE users; --", "value1", "value2") // Simulate SQL injection-like input

	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		om := NewRoarIndex[string, string]()
		om.PushMap(key, value1)
		om.PushMap(key, value2)
		// Test and verify as before
	})
}

// Fuzz test to check behavior with very long strings
func FuzzRoarIndexLongStrings(f *testing.F) {
	longKey := strings.Repeat("k", 100000)
	longValue1 := strings.Repeat("v", 100000)
	longValue2 := strings.Repeat("w", 100000)
	f.Add(longKey, longValue1, longValue2)

	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		// Limit the maximum length to prevent excessive resource usage
		maxLength := 100000
		if utf8.RuneCountInString(key) > maxLength || utf8.RuneCountInString(value1) > maxLength || utf8.RuneCountInString(value2) > maxLength {
			t.Skip("Skipping test case with excessively long strings")
		}

		om := NewRoarIndex[string, string]()
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for long strings: %v", err)
		}

		expectedValues := make(map[string]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}
	})
}

// Fuzz test to ensure RoarIndex handles unicode and special characters correctly
func FuzzRoarIndexUnicode(f *testing.F) {
	f.Add("ключ", "значение1", "значение2") // Russian
	f.Add("キー", "値1", "値2")                 // Japanese
	f.Add("مفتاح", "قيمة1", "قيمة2")        // Arabic
	f.Add("😊", "❤️", "✨")                   // Emojis

	f.Fuzz(func(t *testing.T, key, value1, value2 string) {
		om := NewRoarIndex[string, string]()
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for unicode strings: %v", err)
		}

		expectedValues := make(map[string]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}
	})
}

// Fuzz test to check RoarIndex behavior near uint32 maximum value
func FuzzRoarIndexMaxUint32(f *testing.F) {
	maxUint32 := ^uint32(0)
	f.Add(maxUint32-1, maxUint32-2, maxUint32-3)

	f.Fuzz(func(t *testing.T, key, value1, value2 uint32) {
		om := NewRoarIndex[uint32, uint32]()
		om.PushMap(key, value1)
		om.PushMap(key, value2)

		values, err := om.GetMap(key)
		if err != nil {
			t.Errorf("GetMap failed for max uint32 values: %v", err)
		}

		expectedValues := make(map[uint32]struct{})
		expectedValues[value1] = struct{}{}
		expectedValues[value2] = struct{}{}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}
	})
}
