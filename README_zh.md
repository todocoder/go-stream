
**阅读其他语言版本: [【English](README.md) | [中文】](README_zh.md).**


# Go-Stream 受到 Java 8 Streams 和 go-zero stream的启发

## 简介
&emsp;&emsp;在JAVA中，涉及到对数组、Collection等集合类中的元素进行操作的时候，通常会通过循环的方式进行逐个处理，或者使用Stream的方式进行处理。那么在Go中用的多的是切片，那么这里基于Java的stream的操作习惯用Go语言(
1.18+)的泛型和通道实现了一些简单的流操作功能。

**用Go-Stream 转换或者Groupby后的数据 ，可以直接使用，无需断言。**

> go-stream代码地址：https://github.com/todocoder/go-stream

## Stream 介绍

&emsp;&emsp;可以将Stream流操作分为3种类型：Stream的生成，Stream中间处理，Stream的终止

### Stream的生成

&emsp;&emsp;主要负责新建一个Stream流，或者基于现有的数组创建出新的Stream流。

| API              | 功能说明                                                          |
|------------------|---------------------------------------------------------------|
| Of()             | 通过可变参数`(values ...T)`创建出一个新的stream串行流对象                       |
| OfParallel()     | 通过可变参数`(values ...T)`创建出一个可并行执行stream串行流对象                    |
| OfFrom()         | 通过方法生成`(generate func(source chan<- T))`创建出一个新的stream串行流对象    |
| OfFromParallel() | 通过方法生成`(generate func(source chan<- T))`创建出一个可并行执行stream串行流对象 |
| Concat()         | 多个流拼接的方式创建出一个串行执行stream串行流对象                                  |

### Stream中间处理

&emsp;&emsp;主要负责对Stream进行处理操作，并返回一个新的Stream对象，中间处理操作可以进行叠加。

| API        | 功能说明                                                              |
|------------|-------------------------------------------------------------------|
| Filter()   | 按照条件过滤符合要求的元素， 返回新的stream流                                        |
| Map()      | 按照条件将已有元素转换为另一个对象类型，一对一逻辑，返回新类型的stream流                           |
| FlatMap()  | 按照条件将已有元素转换为另一个对象类型，一对多逻辑，即原来一个元素对象可能会转换为1个或者多个新类型的元素，返回新的stream流 |
| Skip()     | 跳过当前流前面指定个数的元素                                                    |
| Limit()    | 仅保留当前流前面指定个数的元素，返回新的stream流                                       |
| Concat()   | 多个流拼接到当前流下                                                        |
| Distinct() | 按照条件去重符合要求的元素， 返回新的stream流                                        |
| Sorted()   | 按照条件对元素进行排序， 返回新的stream流                                          |
| Reverse()  | 对流中元素进行返转操作                                                       |
| Peek()     | 对stream流中的每个元素进行逐个遍历处理，返回处理后的stream流                              |

### Stream的终止

&emsp;&emsp;通过终止函数操作之后，Stream流将会结束，最后可能会执行某些逻辑处理，或者是按照要求返回某些执行后的结果数据。

| API         | 功能说明                                  |
|-------------|---------------------------------------|
| FindFirst() | 获取第一个元素                               |
| FindLast()  | 获取最后一个元素                              |
| ForEach()   | 对元素进行逐个遍历，然后执行给定的处理逻辑                 |
| Reduce()    | 对流中元素进行聚合处理                           |
| AnyMatch()  | 返回此流中是否存在元素满足所提供的条件                   |
| AllMatch()  | 返回此流中是否全都满足条件                         |
| NoneMatch() | 返回此流中是否全都不满足条件                        |
| Count()     | 返回此流中元素的个数                            |
| Max()       | 返回stream处理后的元素最大值                     |
| Min()       | 返回stream处理后的元素最小值                     |
| ToSlice()   | 将流处理后转化为切片                            |
| Collect()   | 将流转换为指定的类型，通过collectors.Collector进行指定 |

### 转换函数

&emsp;&emsp; 通过这几个函数你可以实现类型转换，分组，flatmap 等处理

> 注意：这几个**函数**非常有用，也是最常用的，由于Go语言泛型的局限性，**Go语言方法**不支持自己独立的泛型，所以导致用Stream中的方法转换只能用 interface{} 代替，这样会有个非常麻烦的问题就是，转换后用的时候必须得强转才能用，所以我把这些写成转换函数，就不会受制于类(struct) 的泛型了。

| API          | 功能说明                                                     |
| ------------ | ------------------------------------------------------------ |
| Map()        | 类型转换(优点：和上面的Map不一样的是，这里转换后可以直接使用，不需要强转) |
| FlatMap()    | 按照条件将已有元素转换为另一个对象类型，一对多逻辑，即原来一个元素对象可能会转换为1个或者多个新类型的元素，返回新的stream流(优点：同Map) |
| GroupingBy() | 对元素进行逐个遍历，然后执行给定的处理逻辑                   |
| Collect()    | 将流转换为指定的类型，通过collectors.Collector进行指定(优点：转换后的类型可以直接使用，无需强转) |

## go-stream的使用

### 库的引入

&emsp;&emsp;由于用到了泛型，支持的版本为golang 1.18+
1. go.mod 中加入如下配置

```
require github.com/todocoder/go-stream v1.1.0
```

2. 执行

```shell
go mod tidy -go=1.20
go mod download
```

### ForEach、Peek的使用

&emsp;&emsp;ForEach和Peek都可以用于对元素进行遍历然后逐个的进行处理。
但Peek属于中间方法，而ForEach属于终止方法。也就是说Peek只能在流中间处理元素，没法直接执行得到结果，其后面必须还要有其它终止操作的时候才会被执行；而ForEach作为无返回值的终止方法，则可以直接执行相关操作。

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

结果如下：

```
item1
item2
item3
item1peek
item2peek
item3peek
```

从代码及结果中得知，ForEach只是用来循环流中的元素。而Peek可以在流中间修改流中的元素。

### Filter、Sorted、Distinct、Skip、Limit、Reverse

&emsp;&emsp;这几个是go-stream中比较常用的中间处理方法，具体说明在上面已标出。使用的话我们可以在流中一个或多个的组合便用。

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
		// 过滤掉1的值
		return item.itemNum != 4
	}).Distinct(func(item TestItem) any {
		// 按itemNum 去重
		return item.itemNum
	}).Sorted(func(a, b TestItem) bool {
		// 按itemNum升序排序
		return a.itemNum < b.itemNum
	}).Skip(1).Limit(6).Reverse().Collect(collectors.ToSlice[TestItem]())
	fmt.Println(res)
}
```

1. 使用Filter过滤掉1的值
2. 通过Distinct对itemNum 去重(在第1步的基础上，下面同理在上一步的基础上)
3. 通过Sorted 按itemNum升序排序
4. 用Skip 从下标为1的元素开始
5. 使用Limit截取排在前6位的元素
6. 使用Reverse 对流中元素进行返转操作
7. 使用collect终止操作将最终处理后的数据收集到Slice中

结果：

```
[{8 item8} {7 item7} {6 item6} {5 item5} {3 item3} {2 item2}]
```

### AllMatch、AnyMatch、NoneMatch、Count、FindFirst、FindLast

&emsp;&emsp;这些方法，均属于这里说的简单结果终止方法。代码如下：

```go
package todocoder

func TestSimple(t *testing.T) {
	allMatch := stream.Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 1, itemValue: "item1"},
	).AllMatch(func(item TestItem) bool {
		// 返回此流中是否全都==1
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
		// 返回此流中是否存在 == 8的
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
		// 返回此流中是否全部不等于8
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

结果：

```
false
true
false
{1 item1} true
{3 item3} true
```

### Map、FlatMap

Map与FlatMap都是用于转换已有的元素为其它元素，区别点在于：

1. Map 按照条件将已有元素转换为另一个对象类型，一对一逻辑
2. FlatMap 按照条件将已有元素转换为另一个对象类型，一对多逻辑

比如我要把 int 1 转为 TestItem{itemNum: 1, itemValue: "item1"}

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

那如果我要把两个字符串["wo shi todocoder","ha ha ha"] 转为 ["wo","shi","todocoder","ha","ha","ha"]
用Map就不行了,这就需要用到FlatMap了

```go
package todocoder

func TestFlatMap(t *testing.T) {
	// 把两个字符串["wo shi todocoder","ha ha ha"] 转为 ["wo","shi","todocoder","ha","ha","ha"]
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

> **注意：这里需要补充一句，只要经过Map或者FlatMap 处理后，类型就会统一变成 `any`了，而不是 泛型`T`，如需要强制类型处理，需要手动转换一下  
> 这个原因是Go泛型的局限性导致的，不能在struct 方法中定义其他类型的泛型，这块看后续官方是否支持了**  
> 
> 可以看如下代码

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
		// 这里需要类型转换
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

&emsp;&emsp;这两个是相对复杂的终止方法，`ToMap` 是类似于Java stream流中`Collectors.toMap()`可以把切片数组转换成 切片map, `GroupBy`
类似于Java stream中 `Collectors.groupingby()`方法，按某个维度来分组

我有如下切片列表：

*TestItem{itemNum: 1, itemValue: "item1"},*  
*TestItem{itemNum: 2, itemValue: "item2"},*  
*TestItem{itemNum: 2, itemValue: "item3"}*  

1. 第一个需求是：把这个列表按 itemNum为Key, itemValue 为 value转换成Map  
2. 第二个需求是：把这个列表按 itemNum为Key, 分组后转换成Map  
我们看一下代码：

```go
package todocoder

func TestToMap(t *testing.T) {
	// 第一个需求
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
	fmt.Println("第一个需求:")
	fmt.Println(resMap)
	// 第二个需求
	resGroup := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item3"},
	).Collect(collectors.GroupingBy(func(t TestItem) any {
		return t.itemNum
	}, func(t TestItem) any {
		return t
	}))
	fmt.Println("第二个需求:")
	fmt.Println(resGroup)
}
```

```
第一个需求:
map[1:item1 2:item2]
第二个需求:
map[1:[{1 item1}] 2:[{2 item2} {2 item3}]]
```

### 转换函数

#### Map、FlatMap(无需断言)

Map与FlatMap都是用于转换已有的元素为其它元素，区别点在于：

1. Map 按照条件将已有元素转换为另一个对象类型，一对一逻辑
2. FlatMap 按照条件将已有元素转换为另一个对象类型，一对多逻辑

比如我要把 TestItem{itemNum: 1, itemValue: "item1"} 转换为ToTestItem{itemNum: 1, itemValue: "item1"}，并且把一个元素按照一定的规则扩展成两个元素，可以通过如下的代码来实现

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

#### GroupingBy() (使用结果无需断言)
需求：
班级有一组学号{1,2,3,....,12}，对应12个人的信息在内存里面存着，把这学号转换成具体的 Student 类，过滤掉 Score 为 1的，并且按评分 Score 分组，并且对各组按照 Age 降序排列

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
    // 注意 这里的返回类型可以是目标类型了
    return studentMap[n]
}).Filter(func(s Student) bool {
    // 这里过滤也不需要转换类型
    // 过滤掉1的
    return s.Score != 1
}), func(t Student) int {
    return t.Score
}, func(t Student) Student {
    return t
}, func(t1 []Student) {
    // 按年龄降序排列
    sort.Slice(t1, func(i, j int) bool {
        return t1[i].Age > t1[j].Age
    })
})
println(res)
```

## 最后

&emsp;&emsp;作为一个Java开发，用习惯了Stream操作，也没找到合适的轻量的stream框架，也不知道后续官方是否会出，在这之前，就先自己简单实现一个，后面遇到复杂的处理流程会持续的更新到上面
除了上面这些功能，还有并行流处理，有兴趣可以自行查看体验[测试类:stream_test](https://github.com/todocoder/go-stream/blob/master/stream/stream_test.go) 

有什么问题可以留言，看到后第一时间回复
