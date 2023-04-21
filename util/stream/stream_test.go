package stream

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zerune/go-core/util/stream/collectors"
)

type TestItem struct {
	itemNum   int
	itemValue string
}

func TestFlatMap2(t *testing.T) {
	// 把两个字符串["wo shi todocoder","ha ha ha"] 转为 ["wo","shi","todocoder","ha","ha","ha"]
	res := Of([]string{"wo shi todocoder", "ha ha ha"}...).FlatMap(func(s string) Stream[any] {
		return OfFrom(func(source chan<- any) {
			for _, str := range strings.Split(s, " ") {
				source <- str
			}
		})
	}).Collect(collectors.ToSlice[any]())
	fmt.Println(res)
}

func TestMap2(t *testing.T) {
	res := Of([]int{1, 2, 3, 4, 7}...).Map(func(item int) any {
		return TestItem{
			itemNum:   item,
			itemValue: fmt.Sprintf("item%d", item),
		}
	}).Collect(collectors.ToSlice[any]())
	fmt.Println(res)
}

type ToTestItem struct {
	itemNum   int
	itemValue string
}

func TestOfFrom(t *testing.T) {
	OfFrom(func(source chan<- TestItem) {
		source <- TestItem{itemNum: 1, itemValue: "item1"}
		source <- TestItem{itemNum: 2, itemValue: "item2"}
		source <- TestItem{itemNum: 3, itemValue: "item3"}
	}).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})
	OfFromParallel(func(source chan<- TestItem) {
		source <- TestItem{itemNum: 1, itemValue: "item1"}
		source <- TestItem{itemNum: 2, itemValue: "item2"}
		source <- TestItem{itemNum: 3, itemValue: "item3"}
	}).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})
}

func TestForEach(t *testing.T) {
	Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})

	OfParallel(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})
}

func TestForEachAndPeek(t *testing.T) {
	// ForEach
	fmt.Println("---------------------ForEach---------------------")
	Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).ForEach(func(item TestItem) {
		fmt.Println(item.itemValue)
	})
	// Peek
	fmt.Println("---------------------Peek---------------------")
	Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Peek(func(item *TestItem) {
		item.itemValue = item.itemValue + "peek"
	}).ForEach(func(item TestItem) {
		fmt.Println(item.itemValue)
	})
}

func TestPeek(t *testing.T) {
	Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Peek(func(item *TestItem) {
		item.itemValue = item.itemValue + "peek"
	}).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})

	OfParallel(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Peek(func(item *TestItem) {
		item.itemValue = item.itemValue + "peek"
	}).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})
}

func TestCount(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Count()
	fmt.Println(res)

	resP := OfParallel(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Count()
	fmt.Println(resP)
}

func TestMaxMin(t *testing.T) {
	resFirst := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindFirst()
	fmt.Println(resFirst.Get())

	resLast := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindLast()
	fmt.Println(resLast.Get())

	resMax := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Max(func(item TestItem, item2 TestItem) int {
		if item.itemNum > item2.itemNum {
			return 1
		} else if item.itemNum == item2.itemNum {
			return 0
		}
		return -1
	})
	fmt.Println(resMax.Get())

	resMin := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Min(func(item TestItem, item2 TestItem) int {
		if item.itemNum > item2.itemNum {
			return 1
		} else if item.itemNum == item2.itemNum {
			return 0
		}
		return -1
	})
	fmt.Println(resMin.Get())
	temp := make([]TestItem, 0)
	resTemp := Of(temp...).Min(func(item TestItem, item2 TestItem) int {
		if item.itemNum > item2.itemNum {
			return 1
		} else if item.itemNum == item2.itemNum {
			return 0
		}
		return -1
	})
	fmt.Println(resTemp.Get())
}

func TestSorted(t *testing.T) {
	resSorted := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Sorted(func(a, b TestItem) int {
		// 降序
		return a.itemNum - b.itemNum
	}).ToSlice()
	fmt.Println(resSorted)
}

func TestReverse(t *testing.T) {
	resReverse := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Reverse().ToSlice()
	fmt.Println(resReverse)
}

func TestDistinct(t *testing.T) {
	resReverse := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 3, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item4"},
	).Distinct(func(item TestItem) any {
		//return item.itemNum
		return item.itemValue
	}).ToSlice()
	fmt.Println(resReverse)
}

func TestSkip(t *testing.T) {
	resSkip := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 3, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item4"},
	).Skip(1).ToSlice()
	fmt.Println(resSkip)
}

func TestLimit(t *testing.T) {
	resLimit := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 3, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item4"},
	).Skip(1).Limit(3).ToSlice()
	fmt.Println(resLimit)
}

func TestReduce(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Reduce(func(item TestItem, item2 TestItem) TestItem {
		return TestItem{itemNum: item.itemNum + item2.itemNum, itemValue: fmt.Sprintf("%s_%s", item.itemValue, item2.itemValue)}
	})
	fmt.Println(res.Get())
}

func TestToSlice(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.ToSlice[TestItem]()).([]TestItem)
	fmt.Println(res)
	//sum := func(i1, i2 int) int { return i1 + i2 }
}

func TestToMap(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item5"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.ToMap[TestItem](func(t TestItem) any {
		return t.itemNum
	}, func(item TestItem) any {
		return item
	}))
	fmt.Println(res)
}

func TestToMapOpt(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item5"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.ToMap[TestItem](func(t TestItem) any {
		return t.itemNum
	}, func(item TestItem) any {
		return item.itemValue
	}, func(oldV, newV any) any {
		return oldV
	}))
	fmt.Println(res)
}

func TestToGroupBy(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item5"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.GroupingBy(func(t TestItem) any {
		return t.itemNum
	}, func(t TestItem) any {
		return t
	}))
	fmt.Println(res)
}

func TestFilter(t *testing.T) {
	Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 9, itemValue: "item9"},
	).Filter(func(item TestItem) bool {
		if item.itemNum%2 != 0 {
			return true
		}
		return false
	}).ForEach(func(item TestItem) {
		fmt.Print(item.itemNum)
		fmt.Println(item.itemValue)
	})
}

func TestMap(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 9, itemValue: "item9"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Map(func(item TestItem) any {
		return ToTestItem{
			itemNum:   item.itemNum,
			itemValue: item.itemValue,
		}
	}).Collect(collectors.ToSlice[any]()).([]any)
	for _, r := range res {
		fmt.Println(r)
	}
}

func TestMapToInt(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 9, itemValue: "item9"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).MapToInt(func(item TestItem) int {
		return item.itemNum
	}).Collect(collectors.ToSlice[int]()).([]int)
	fmt.Println(res)
}

func TestFlatMap(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).FlatMap(func(item TestItem) Stream[any] {
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

func TestConcatenate(t *testing.T) {
	resConcat := Concat(Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
	), Of(
		TestItem{itemNum: 5, itemValue: "item5"},
		TestItem{itemNum: 6, itemValue: "item6"},
	), Of(
		TestItem{itemNum: 3, itemValue: "item3"},
		TestItem{itemNum: 4, itemValue: "item4"},
	)).Collect(collectors.ToSlice[TestItem]())
	fmt.Println(resConcat)
}

func TestSimple(t *testing.T) {
	allMatch := Of(
		TestItem{itemNum: 7, itemValue: "item7"},
		TestItem{itemNum: 6, itemValue: "item6"},
		TestItem{itemNum: 8, itemValue: "item8"},
		TestItem{itemNum: 1, itemValue: "item1"},
	).AllMatch(func(item TestItem) bool {
		// 返回此流中是否全都==1
		return item.itemNum == 1
	})
	fmt.Println(allMatch)

	anyMatch := Of(
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

	noneMatch := Of(
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

	resFirst := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindFirst()
	fmt.Println(resFirst.Get())

	resLast := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).FindLast()
	fmt.Println(resLast.Get())
}

func TestToMap11(t *testing.T) {
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
	fmt.Println(resMap)

	resGroup := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item3"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.GroupingBy(func(t TestItem) any {
		return t.itemNum
	}, func(t TestItem) any {
		return t
	}))
	fmt.Println(resGroup)
}

func TestToMapOpt11(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item2"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).Collect(collectors.ToMap[TestItem](func(t TestItem) any {
		return t.itemNum
	}, func(item TestItem) any {
		return item.itemValue
	}, func(oldV, newV any) any {
		return oldV
	}))
	fmt.Println(res)
}
