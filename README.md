&emsp;&emsp;

**Read this in other languages: [【English](README.md) | [中文】](README_zh.md).**

# Stream Collections for Go. Inspired in Java 8 Streams and go-zero.

## 简介
&emsp;&emsp;In JAVA, when it comes to operating elements in collection classes such as arrays and Collections, they are usually processed one by one through a loop, or processed using a Stream. Well, slices are mostly used in Go, so the Java-based stream operations here are customary to use the Go language (
1.18+)'s generics and channels implement some simple stream operation functions.

**Data converted with Go-Stream or Groupby can be used directly without assertions.**

> go-stream code address：https://github.com/todocoder/go-stream

## Stream Introduction

&emsp;&emsp;Stream operations can be divided into three types: Stream generation, Stream intermediate processing, and Stream termination.

### Stream generation

&emsp;&emsp;Mainly responsible for creating a new Stream or creating a new Stream based on an existing array.

| API              | Function Description                                                                                                   |
|------------------|------------------------------------------------------------------------------------------------------------------------|
| Of()             | Create a new stream serial stream object through variable parameters `(values ...T)`                                   |
| OfParallel()     | Create a stream serial stream object that can be executed in parallel through variable parameters `(values ...T)`      |
| OfFrom()         | Create a new stream serial stream object through the method `(generate func(source chan<- T))`                         |
| OfFromParallel() | Generate a serial stream object that can be executed in parallel through the method `(generate func(source chan<- T))` |
| Concat()         | Multiple streams are spliced together to create a serial execution stream serial stream object.                        |

### Stream intermediate processing

&emsp;&emsp;Mainly responsible for processing Stream and returning a new Stream object. Intermediate processing operations can be superimposed.

| API        | Function Description                                                                                                                                                                                                      |
|------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Filter()   | Filter elements that meet the requirements according to conditions and return a new stream                                                                                                                                |
| Map()      | Convert existing elements to another object type according to conditions, one-to-one logic, and return a new type of stream                                                                                               |
| FlatMap()  | Convert existing elements to another object type according to conditions, one-to-many logic, that is, the original element object may be converted into one or more elements of a new type, and a new stream is returned. |
| Skip()     | Skip the specified number of elements in front of the current stream                                                                                                                                                      |
| Limit()    | Only retain the specified number of elements in front of the current stream and return a new stream                                                                                                                       |
| Concat()   | Multiple streams are spliced under the current stream                                                                                                                                                                     |
| Distinct() | Eliminate duplicate elements that meet the requirements according to conditions and return a new stream.                                                                                                                  |
| Sorted()   | Sort elements according to conditions and return a new stream                                                                                                                                                             |
| Reverse()  | Reverse elements in a stream                                                                                                                                                                                              |
| Peek()     | Traverse each element in the stream one by one and return the processed stream                                                                                                                                            |

### Stream termination

&emsp;&emsp;After the termination function operation, the Stream will end, and finally some logical processing may be performed, or some execution result data may be returned as required.

| API         | Function Description                                                                     |
|-------------|------------------------------------------------------------------------------------------|
| FindFirst() | Get the first element                                                                    |
| FindLast()  | Get the last element                                                                     |
| ForEach()   | Traverse the elements one by one and then execute the given processing logic             |
| Reduce()    | Aggregate elements in a stream                                                           |
| AnyMatch()  | Returns whether there is an element in this stream that satisfies the provided condition |
| AllMatch()  | Returns whether all conditions in this stream are met                                    |
| NoneMatch() | Returns whether all conditions in this stream are not met                                |
| Count()     | Returns the number of elements in this stream                                            |
| Max()       | Returns the maximum value of the element after stream processing                         |
| Min()       | Returns the minimum value of the element after stream processing                         |
| ToSlice()   | Convert streams into slices after processing                                             |
| Collect()   | Convert the stream to the specified type, specified through collectors.Collector         |

### Conversion Function

&emsp;&emsp;Through these functions you can implement type conversion, grouping, flatmap and other processing

> Note: These **functions** are very useful and the most commonly used. Due to the limitations of Go language generics, **Go language methods** do not support their own independent generics, so the method in Stream is used for conversion. It can only be replaced by interface{}. This will have a very troublesome problem. It must be forced to be used after conversion, so I wrote these as conversion functions so that they will not be subject to the generics of the class (struct).。

| API          | Function Description                                                                                                                                                                                                                             |
|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Map()        | Type conversion (advantage: unlike the Map above, it can be used directly after conversion without forced conversion)                                                                                                                            |
| FlatMap()    | Convert existing elements to another object type according to conditions, one-to-many logic, that is, an original element object may be converted into one or more elements of a new type, and a new stream is returned (advantage: same as Map) |
| GroupingBy() | Traverse the elements one by one and then execute the given processing logic                                                                                                                                                                     |
| Collect()    | Convert the stream to the specified type and specify it through collectors.Collector (advantage: the converted type can be used directly without forced conversion)                                                                              |

## Use of Go-Stream

### Introduce

&emsp;&emsp;Due to the use of generics, the supported version is golang 1.18+
1. Add the following configuration to go.mod

```
require github.com/todocoder/go-stream v1.1.0
```

2. Implement

```shell
go mod tidy -go=1.20
go mod download
```

### Use of ForEach and Peek

&emsp;&emsp;Both ForEach and Peek can be used to traverse elements and process them one by one.
But Peek is an intermediate method, while ForEach is a termination method. In other words, Peek can only process elements in the middle of the stream, and cannot directly execute it to obtain the result. It will only be executed when there are other termination operations; and ForEach, as a termination method without a return value, can directly execute the relevant operate.

#### ForEach

```go
package todocoder

type TestItem struct {
	itemNum   int
	itemValue string
}

func TestForEachAndPeek(t *testing.T) {
	// ForEach
	stream.Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).ForEach(func(item TestItem) {
		fmt.Println(item.itemValue)
	})
	// Peek
	stream.Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Peek(func(item *TestItem) {
		item.itemValue = item.itemValue + "peek"
	}).ForEach(func(item TestItem) {
		fmt.Println(item.itemValue)
	})
}
```

The result is as follows:

```
item1
item2
item3
item1peek
item2peek
item3peek
```

From the code and results, we know that ForEach is only used to loop through the elements in the stream. Peek can modify elements in the stream in the middle of the stream.

### Filter、Sorted、Distinct、Skip、Limit、Reverse

&emsp;&emsp;These are the more commonly used intermediate processing methods in go-stream, and the specific instructions are marked above. If used, we can use one or more combinations in the stream.

```go
package todocoder

func TestStream(t *testing.T) {
	// ForEach
	res := stream.Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 9, itemValue: "item9"},
	).Filter(func(item TestItem) bool {
		// Filter out values of 4
		return item.itemNum != 4
	}).Distinct(func(item TestItem) any {
		// Press itemNum to remove duplicates
		return item.itemNum
	}).Sorted(func(a, b TestItem) bool {
		// Sort by itemNum in ascending order
		return a.itemNum < b.itemNum
	}).Skip(1).Limit(6).Reverse().Collect(collectors.ToSlice[TestItem]())
	fmt.Println(res)
}
```

1. Use Filter() to filter out the value of 4
2. Deduplicate itemNum through Distinct() (based on the first step, the following is the same as the previous step)
3. Sorted() by itemNum in ascending order by Sorted
4. Use Skip() to start from the element with index 1
5. Use Limit() to intercept the top 6 elements
6. Use Reverse() to reverse elements in a stream
7. Use Collect() to terminate the operation and collect the final processed data into Slice

result：

```
[{8 item8} {7 item7} {6 item6} {5 item5} {3 item3} {2 item2}]
```

### AllMatch、AnyMatch、NoneMatch、Count、FindFirst、FindLast

&emsp;&emsp;These methods all belong to the simple result termination methods mentioned here. code show as below:

```go
package todocoder

func TestSimple(t *testing.T) {
	allMatch := stream.Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 1, itemValue: "item1"},
	).AllMatch(func(item TestItem) bool {
		// Returns whether all in this stream == 1
		return item.itemNum == 1
	})
	fmt.Println(allMatch)

	anyMatch := stream.Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 1, itemValue: "item1"},
	).Filter(func(item TestItem) bool {
		return item.itemNum != 1
	}).AnyMatch(func(item TestItem) bool {
		// Returns whether there is == 8 in this stream
		return item.itemNum == 8
	})
	fmt.Println(anyMatch)

	noneMatch := stream.Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 1, itemValue: "item1"},
	).Filter(func(item TestItem) bool {
		return item.itemNum != 1
	}).NoneMatch(func(item TestItem) bool {
		// Returns whether all in this stream are not equal to 8
		return item.itemNum == 8
	})
	fmt.Println(noneMatch)

	resFirst := stream.Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindFirst()
	fmt.Println(resFirst.Get())

	resLast := stream.Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindLast()
	fmt.Println(resLast.Get())
}
```

result：

```
false
true
false
{1 item1} true
{3 item3} true
```

### Map、FlatMap

Both Map and FlatMap are used to convert existing elements into other elements. The difference is:

1. Map() converts existing elements into another object type according to conditions, one-to-one logic
2. FlatMap() converts existing elements to another object type according to conditions, one-to-many logic

For example, I want to convert int 1 to TestItem{itemNum: 1, itemValue: "item1"}

```go
package todocoder

func TestMap(t *testing.T) {
	res := stream.Of([]int{1, 2, 3, 4, 7}...).Map(func(item int) any {
		return TestItem{
			itemNum:   item,
			itemValue: fmt.Sprintf("item%d", item),
		}
	}).Collect(collectors.ToSlice[any]())
	fmt.Println(res)
}
```

```
[{1 item1} {2 item2} {3 item3} {4 item4} {7 item7}]
```

So if I want to convert two strings ["wo shi todocoder", "ha ha ha"] into ["wo", "shi", "todocoder", "ha", "ha", "ha"]
Using Map won’t work, so you need to use FlatMap

```go
package todocoder

func TestFlatMap(t *testing.T) {
	// Convert two strings ["wo shi todocoder", "ha ha ha"] into ["wo", "shi", "todocoder", "ha", "ha", "ha"]
	res := stream.Of([]string{"wo shi todocoder", "ha ha ha"}...).FlatMap(func(s string) stream.Stream[any] {
		return stream.OfFrom(func(source chan<- any) {
			for _, str := range strings.Split(s, " ") {
				source <- str
			}
		})
	}).Collect(collectors.ToSlice[any]())
	fmt.Println(res)
}
```

```
[wo shi todocoder ha ha ha]
```

> **Note: It is necessary to add here that as long as it is processed by Map or FlatMap, the type will become `any` instead of generic `T`. If forced type processing is required, manual conversion is required.
> This reason is caused by the limitations of Go generics. Other types of generics cannot be defined in the struct method. This will depend on whether it will be officially supported in the future.**
> 
> You can see the following code:

```go
package todocoder

func TestMap(t *testing.T) {
	res := stream.Of(
		TestItem{itemNum: 3, itemValue: "item3"},
	).FlatMap(func(item TestItem) stream.Stream[any] {
		return Of[any](
			TestItem{itemNum: item.itemNum * 10, itemValue: fmt.Sprintf("%s+%d", item.itemValue, item.itemNum)},
			TestItem{itemNum: item.itemNum * 20, itemValue: fmt.Sprintf("%s+%d", item.itemValue, item.itemNum)},
		)
	}).Map(func(item any) any {
		// Type assertion is required here
		ite := item.(TestItem)
		return ToTestItem{
			itemNum:   ite.itemNum,
			itemValue: ite.itemValue,
		}
	}).Collect(collectors.ToSlice[any]())
	fmt.Println(res)
}
```

### collectors.ToMap、collectors.GroupBy

&emsp;&emsp;These two are relatively complex termination methods. `ToMap` is similar to `Collectors.toMap()` in Java stream, which can convert the slice array into a slice map, `GroupBy`
Similar to the `Collectors.groupingby()` method in Java stream, grouped by a certain dimension

The following slice list：

*TestItem{itemNum: 1, itemValue: "item1"},*  
*TestItem{itemNum: 2, itemValue: "item2"},*  
*TestItem{itemNum: 2, itemValue: "item3"}*  

1. The first requirement is: convert this list into a Map according to itemNum as Key and itemValue as value. 
2. The second requirement is: group this list according to itemNum as Key, group it and convert it into a Map

Let's take a look at the code：

```go
package todocoder

func TestToMap(t *testing.T) {
	// The first requirement
	resMap := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item3"},
	).Collect(collectors.ToMap[TestItem](func(t TestItem) any {
		return t.itemNum
	}, func(item TestItem) any {
		return item.itemValue
	}, func(oldV, newV any) any {
		return oldV
	}))
	fmt.Println("The first requirement:")
	fmt.Println(resMap)
	// The second requirement
	resGroup := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item3"},
	).Collect(collectors.GroupingBy(func(t TestItem) any {
		return t.itemNum
	}, func(t TestItem) any {
		return t
	}))
	fmt.Println("The second requirement:")
	fmt.Println(resGroup)
}
```

```
The first requirement:
map[1:item1 2:item2]
The second requirement:
map[1:[{1 item1}] 2:[{2 item2} {2 item3}]]
```

### Conversion Function

### Map、FlatMap(No assertion required)

Both Map and FlatMap are used to convert existing elements into other elements. The difference is:

1. Map() converts existing elements into another object type according to conditions, one-to-one logic
2. FlatMap() converts existing elements to another object type according to conditions, one-to-many logic

For example, if I want to convert TestItem{itemNum: 1, itemValue: "item1"} to ToTestItem{itemNum: 1, itemValue: "item1"}, and expand one element into two elements according to certain rules, I can use the following code to fulfill
```go
func TestFlatMap(t *testing.T) {
	res := stream.Map(stream.Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FlatMap(func(item TestItem) Stream[TestItem] {
		return Of[TestItem](
			TestItem{itemNum: item.itemNum * 10, itemValue: fmt.Sprintf("%s+%d", item.itemValue, item.itemNum)},
			TestItem{itemNum: item.itemNum * 20, itemValue: fmt.Sprintf("%s+%d", item.itemValue, item.itemNum)},
		)
	}), func(item TestItem) ToTestItem {
		return ToTestItem{
			itemNum:   item.itemNum,
			itemValue: item.itemValue,
		}
	}).ToSlice()
	fmt.Println(res)
}
```

#### GroupingBy() (Use results without assertions)
Need:
The class has a set of student numbers `{1,2,3,....,12}`, and the information corresponding to 12 people is stored in the memory. Convert this student number into a specific Student class, filter out those with a Score of 1, and Group by Score and sort each group in descending order by Age

```go

studentMap := map[int]Student{
    1:  {Num: 1, Name: "小明", Score: 3, Age: 26},
    2:  {Num: 2, Name: "小红", Score: 4, Age: 27},
    3:  {Num: 3, Name: "小李", Score: 5, Age: 19},
    4:  {Num: 4, Name: "老王", Score: 1, Age: 23},
    5:  {Num: 5, Name: "小王", Score: 2, Age: 29},
    6:  {Num: 6, Name: "小绿", Score: 2, Age: 24},
    7:  {Num: 7, Name: "小蓝", Score: 3, Age: 29},
    8:  {Num: 8, Name: "小橙", Score: 3, Age: 30},
    9:  {Num: 9, Name: "小黄", Score: 4, Age: 22},
    10: {Num: 10, Name: "小黑", Score: 5, Age: 21},
    11: {Num: 11, Name: "小紫", Score: 3, Age: 32},
    12: {Num: 12, Name: "小刘", Score: 2, Age: 35},
}

res := GroupingBy(Map(Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12), func(n int) Student {
    // Note that the return type here is the target type, no assertion is needed
    return studentMap[n]
}).Filter(func(s Student) bool {
    // There is no need to convert types when filtering here.
    return s.Score != 1
}), func(t Student) int {
    return t.Score
}, func(t Student) Student {
    return t
}, func(t1 []Student) {
    // Sort by age in descending order
    sort.Slice(t1, func(i, j int) bool {
        return t1[i].Age > t1[j].Age
    })
})
println(res)
```


## At Last

&emsp;&emsp;As a Java developer, I am used to Stream operations, but I haven’t found a suitable lightweight stream framework. I don’t know whether the official one will be released in the future. Before that, I will simply implement one by myself. I will encounter complex processing processes later. Continuously update to the above
In addition to the above functions, there is also parallel stream processing. If you are interested, you can check the experience by yourself [Test category: stream_test](https://github.com/todocoder/go-stream/blob/master/stream/stream_test.go)

If you have any questions, please leave a message and we will reply as soon as possible after seeing it.
