package roarindex_test

import (
	"fmt"
	"slices"

	"roarindex"
)

func ExampleRoarIndex() {
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap", "value2")
	cm.PushMap("testMap", "value3")

	values, _ := cm.GetMap("testMap")
	// order is not guaranteed
	slices.Sort(values)
	fmt.Println(values)
	// Output: [value1 value2 value3]
}

func ExampleRoarIndex_Count() {
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap", "value2")
	cm.PushMap("testMap", "value3")
	fmt.Println(cm.Count())
	// Output: 1
}

func ExampleRoarIndex_Values() {
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap1", "value1")
	cm.PushMap("testMap2", "value2")
	cm.PushMap("testMap3", "value2")

	values := cm.Values()
	// order is not guaranteed
	slices.Sort(values)
	fmt.Println(values)
	// Output: [value1 value2]
}

func ExampleRoarIndex_DeleteMap() {
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.DeleteMap("testMap")
	fmt.Println(cm.GetMap("testMap"))
	// Output: [] key not found
}

func ExampleRoarIndex_HasValue() {
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap2", "value1")

	fmt.Printf("%t %t %t", cm.HasValue("testMap", "value1"), cm.HasValue("testMap2", "value1"), cm.HasValue("testMap", "value2"))
	// Output: true true false
}
