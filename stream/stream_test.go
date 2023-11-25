package stream

import (
	"encoding/json"
	"fmt"
	"github.com/todocoder/go-stream/collectors"
	"sort"
	"strings"
	"testing"
)

type TestItem struct {
	itemNum   int
	itemValue string
}
type TestItemS struct {
	itemNum   int    `json:"itemNum"`
	itemValue string `json:"itemValue"`
}

func TestFlatMap2(t *testing.T) {
	// 把两个字符串["wo shi todocoder","ha ha ha"] 转为 ["wo","shi","todocoder","ha","ha","ha"]
	res := Of([]string{"wo shi todocoder", "ha ha ha"}...).FlatMap(func(s string) Stream[string] {
		return OfFrom(func(source chan<- string) {
			for _, str := range strings.Split(s, " ") {
				source <- str
			}
		})
	}).ToSlice()

	fmt.Println(res)
}

func TestMap2(t *testing.T) {
	res := Of([]int{1, 2, 3, 4, 7}...).Map(func(item int) any {
		return TestItem{
			itemNum:   item,
			itemValue: fmt.Sprintf("item%d", item),
		}
	}).ToSlice()
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

func TestToAverage(t *testing.T) {
	resInt := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).MapToInt(func(item TestItem) int {
		return item.itemNum
	}).SumIntStatistics()
	fmt.Println(resInt.GetSum())
	fmt.Println(resInt.GetCount())
	fmt.Println(resInt.GetMin())
	fmt.Println(resInt.GetMax())
	fmt.Println(resInt.GetAverage())

	resFloat := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).MapToFloat64(func(item TestItem) float64 {
		return float64(item.itemNum)
	}).SumIntStatistics()
	fmt.Println(resFloat.GetSum())
	fmt.Println(resFloat.GetCount())
	fmt.Println(resFloat.GetMin())
	fmt.Println(resFloat.GetMax())
	fmt.Println(resFloat.GetAverage())
}

func TestSorted(t *testing.T) {
	resSorted := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Sorted(func(a, b TestItem) bool {
		// 降序
		return a.itemNum > b.itemNum
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
	).Skip(1).Limit(7).ToSlice()
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
	}).ToSlice()
	fmt.Println(res)
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
	}).ToMapInt(func(t TestItem) int {
		return t.itemNum
	}, func(item TestItem) TestItem {
		return item
	})
	fmt.Println(res)
}

func TestToMap1(t *testing.T) {
	res0 := Collect(Of(
		TestItemS{itemNum: 1, itemValue: "item1"},
		TestItemS{itemNum: 2, itemValue: "item2"},
		TestItemS{itemNum: 3, itemValue: "item3"},
		TestItemS{itemNum: 4, itemValue: "item4"},
		TestItemS{itemNum: 4, itemValue: "item4"},
		TestItemS{itemNum: 4, itemValue: "item5"},
	), collectors.ToMap(func(t TestItemS) string {
		return t.itemValue
	}, func(item TestItemS) int {
		return item.itemNum
	}))
	println(res0)

	res1 := Collect(Of(
		TestItemS{itemNum: 1, itemValue: "item1"},
		TestItemS{itemNum: 2, itemValue: "item2"},
		TestItemS{itemNum: 3, itemValue: "item3"},
		TestItemS{itemNum: 4, itemValue: "item4"},
		TestItemS{itemNum: 5, itemValue: "item4"},
		TestItemS{itemNum: 4, itemValue: "item5"},
	), collectors.ToMap(func(t TestItemS) string {
		return t.itemValue
	}, func(item TestItemS) int {
		return item.itemNum
	}, func(v1, v2 int) int {
		return v2
	}))
	println(res1)
	v, k := json.Marshal(res0)

	fmt.Println(v)
	fmt.Println(k)
}

func TestToMapOpt(t *testing.T) {
	res := Collect(Of(
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
	}), collectors.ToMap(func(t TestItem) int {
		return t.itemNum
	}, func(item TestItem) string {
		return item.itemValue
	}, func(oldV, newV string) string {
		return oldV
	}))
	fmt.Println(res)
}

func TestToGroupBy(t *testing.T) {
	res := Collect(Of(
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
	}), collectors.GroupingBy(func(t TestItem) int {
		return t.itemNum
	}, func(t TestItem) string {
		return t.itemValue
	}))
	fmt.Println(res)

	res1 := Collect(Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item2"},
		TestItem{itemNum: 5, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 0, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item5"},
		TestItem{itemNum: 9, itemValue: "item5"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}), collectors.GroupingBy(func(t TestItem) string {
		return t.itemValue
	}, func(t TestItem) int {
		return t.itemNum
	}, func(t1 []int) {
		sort.Slice(t1, func(i, j int) bool {
			return t1[i] < t1[j]
		})
	}))
	fmt.Println(res1)

	res2 := GroupingBy(Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item2"},
		TestItem{itemNum: 5, itemValue: "item2"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 4, itemValue: "item4"},
		TestItem{itemNum: 0, itemValue: "item4"},
		TestItem{itemNum: 4, itemValue: "item5"},
		TestItem{itemNum: 9, itemValue: "item5"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}), func(t TestItem) string {
		return t.itemValue
	}, func(t TestItem) int {
		return t.itemNum
	}, func(t1 []int) {
		sort.Slice(t1, func(i, j int) bool {
			return t1[i] < t1[j]
		})
	})
	fmt.Println(res2)

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
	}).ToSlice()
	fmt.Println(res)
}

func TestFlatMap(t *testing.T) {
	res := Map(Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).Filter(func(item TestItem) bool {
		if item.itemNum != 1 {
			return true
		}
		return false
	}).FlatMap(func(item TestItem) Stream[TestItem] {
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
	)).ToSlice()
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

func TestGroupby(t *testing.T) {
	res := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "3tem21"},
		TestItem{itemNum: 2, itemValue: "2tem22"},
		TestItem{itemNum: 2, itemValue: "1tem22"},
		TestItem{itemNum: 2, itemValue: "4tem23"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).GroupingByInt(func(t TestItem) int {
		return t.itemNum
	})
	fmt.Println(res)
	res1 := Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 6, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).GroupingByString(func(t TestItem) string {
		return t.itemValue
	}, func(t1 []TestItem) {
		sort.Slice(t1, func(i, j int) bool {
			return t1[i].itemNum < t1[j].itemNum
		})
	})
	fmt.Println(res1)
}

func TestStatistic(t *testing.T) {
	res := Collect(Of(
		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "3tem21"},
		TestItem{itemNum: 2, itemValue: "2tem22"},
		TestItem{itemNum: 2, itemValue: "1tem22"},
		TestItem{itemNum: 2, itemValue: "4tem23"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).MapToInt(func(t TestItem) int {
		return t.itemNum * 2
	}), collectors.Statistic[int]())

	println(res.GetAverage())
	res1 := Of([]int64{1, 2, 6, 8, 9, 20}...).SumInt64Statistics()
	println(res1.GetSum())
	println(res1.GetCount())
	println(res1.GetMax())
	println(res1.GetMin())
	println(res1.GetAverage())

	res2 := Of([]int32{1, 2, 6, 8, 9, 20}...).SumInt32Statistics()
	println(res2.GetSum())

	res3 := Of([]int{1, 2, 6, 8, 9, 20}...).SumIntStatistics()
	println(res3.GetSum())

	res4 := Of([]float32{1, 2, 6, 8, 9, 20}...).SumFloat32Statistics()
	println(res4.GetSum())

	res5 := Of([]float32{1, 2, 6, 8, 9, 20}...).SumFloat64Statistics()
	println(res5.GetSum())
}
