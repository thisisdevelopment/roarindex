# RoarIndex

RoarIndex is a high-performance, memory-efficient indexing data structure for Go that leverages Roaring Bitmaps to provide fast set operations and lookups. It is designed to handle large-scale data indexing with optimal performance.

## Features

- Fast set and get operations
- Memory efficient storage of large sets
- Thread-safe concurrent access
- Generic implementation supporting different key-value type combinations
- Automatic deduplication of values
- Efficient lookup and retrieval operations

## Why RoarIndex?

RoarIndex was created to solve the challenge of efficiently indexing and querying large sets of data while maintaining good memory usage characteristics. It uses Roaring Bitmaps under the hood, which are a compressed bitmap data structure that has been shown to be significantly more efficient than traditional bitmap implementations.

### What are Roaring Bitmaps?

[Roaring Bitmaps](https://roaringbitmap.org/) are a hybrid data structure that combines different compression techniques to efficiently store and process sets of integers. They are particularly effective when dealing with sparse data sets and have been adopted by several big data systems including Apache Lucene, Apache Spark, and Apache Druid.


## Benchmarks
```bash
❯ go test -benchmem -bench=.
goos: darwin
goarch: arm64
pkg: roarindex
cpu: Apple M2 Max
BenchmarkRoarIndexPushMap-12            10301625               115.4 ns/op            13 B/op          1 allocs/op
BenchmarkRoarIndexGetMap-12              5075074               236.0 ns/op           205 B/op          3 allocs/op
BenchmarkRoarIndexHasValue-12           10324004               117.3 ns/op            13 B/op          1 allocs/op
BenchmarkRoarIndexDeleteMap-12          21029878                57.40 ns/op           13 B/op          1 allocs/op
BenchmarkMemoryUsageComparison/RoarIndex-500000000-kv-12                       1        22166860000 ns/op               87.62 MB-RoarIndex      339899048 B/op    224019 allocs/op
BenchmarkMemoryUsageComparison/StandardMap-500000000-kv-12                     1        1442234500 ns/op              4529 MB-StdMap    19479252520 B/op          240145 allocs/op
PASS
ok      roarindex       34.278s
```


### Memory Usage Comparison

The benchmark results show a dramatic difference in memory usage between RoarIndex and a standard Go map:

- RoarIndex: **~87.62** MB memory usage
- Standard Map: ~4,529 MB memory usage

This represents a **98% reduction in memory usage** when using RoarIndex compared to a standard map for the same dataset of 500 million key-value pairs.

The key factors behind this massive memory efficiency:

1. **Bitmap Compression**: RoarIndex uses Roaring Bitmaps which employ sophisticated compression techniques to store integer sets very efficiently

2. **Deduplication**: RoarIndex automatically deduplicates values, storing each unique value only once

3. **Optimized Storage**: The bitmap-based storage eliminates the overhead of storing duplicate pointers and metadata that a standard map requires for each entry

4. **Memory Allocation Efficiency**: RoarIndex makes significantly fewer allocations (224,019 vs 240,145) and allocates much less memory per operation (339MB vs 19,479MB total)

This memory efficiency makes RoarIndex particularly well-suited for large-scale applications where memory usage is a critical concern, such as in-memory databases, caching systems, and high-performance data processing pipelines.

## Usage

```go
	cm := roarindex.NewRoarIndex[string, string]()
```
initializes a new RoarIndex with string keys and values, in this case a **[]string** for the value.

```go
	cm := roarindex.NewRoarIndex[uint64, SomeTypeStruct]()
```
initializes a new RoarIndex with uint64 keys and SomeTypeStruct values, in this case a **[]SomeTypeStruct** for the value. 

pushing a value to the map:

```go
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap", "value2")
	cm.PushMap("testMap", "value3")
	fmt.Println(cm.GetMap("testMap"))
	// Output: [value1 value2 value3] <nil>
```

```go
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap", "value2")
	cm.PushMap("testMap", "value3")
	fmt.Println(cm.Count())
	// Output: 1
```
returns the number of keys in the RoarIndex.

```go
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap1", "value1")
	cm.PushMap("testMap2", "value2")
	cm.PushMap("testMap3", "value3")
	fmt.Println(cm.Values())
	// Output: [value1 value2 value3]
```
returns all the values in the RoarIndex.

```go
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.DeleteMap("testMap")
	fmt.Println(cm.GetMap("testMap"))
	// Output: [] key not found
```
deletes the "testMap" key and all values associated with it.

```go
	cm := roarindex.NewRoarIndex[string, string]()
	cm.PushMap("testMap", "value1")
	cm.PushMap("testMap2", "value1")

	fmt.Printf("%t %t %t", cm.HasValue("testMap", "value1"), cm.HasValue("testMap2", "value1"), cm.HasValue("testMap", "value2"))
	// Output: true true false
```
checks if the value exists in the "testMap" key.

## About Us Th[is]

[This](https://this.nl) is a digital agency based in Utrecht, the Netherlands, specializing in crafting high-performance, resilient, and scalable digital solutions, api's, microservices, and more. Our multidisciplinary team of designers, front and backend developers and strategists collaborates closely to deliver robust and efficient products that meet the demands of today's digital landscape. We are passionate about turning ideas into reality and providing exceptional value to our clients through innovative technology and exceptional user experiences.

## Contributing

Contributions are welcome! We especially encourage contributions of new storage backends. Please open an issue to discuss your ideas or submit a pull request with your implementation.

## License

This project is licensed under the MIT License.
