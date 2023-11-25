package collectors

import (
	"github.com/todocoder/go-stream/utils"
)

/*
Number 抽象出常用的运算类型
*/
type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

/*
SumStatistics 这个struct 是一些简单统计的聚合类
*/
type SumStatistics[S Number] struct {
	Count   int64
	Sum     S
	Max     S
	Min     S
	Average float64
}

func (s SumStatistics[S]) GetCount() int64 {
	return s.Count
}
func (s SumStatistics[S]) GetSum() S {
	return s.Sum
}
func (s SumStatistics[S]) GetMax() S {
	return s.Max
}
func (s SumStatistics[S]) GetMin() S {
	return s.Min
}
func (s SumStatistics[S]) GetAverage() float64 {
	return s.Average
}

/*
Statistic Collect 统计方法
可以和stream结合使用，如下的代码来使用：

eg:

	type TestItem struct {
		itemNum   int
		itemValue string
	}

res := stream.Collect(stream.Of(

		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "3tem21"},
		TestItem{itemNum: 2, itemValue: "2tem22"},
		TestItem{itemNum: 2, itemValue: "1tem22"},
		TestItem{itemNum: 2, itemValue: "4tem23"},
		TestItem{itemNum: 3, itemValue: "item3"},
	).MapToInt(func(t TestItem) int {
		return t.itemNum * 2
	}), collectors.Statistic[int]())
	println(res.GetSum())
	println(res.GetAverage())
*/
func Statistic[T Number]() Collector[T, *SumStatistics[T], SumStatistics[T]] {
	return &DefaultCollector[T, *SumStatistics[T], SumStatistics[T]]{
		supplier: func() *SumStatistics[T] {
			return &SumStatistics[T]{
				Count:   0,
				Sum:     0,
				Max:     0,
				Min:     0,
				Average: 0,
			}
		},
		accumulator: func(ts *SumStatistics[T], t T) {
			if ts.Count == 0 {
				ts.Min = t
			}
			ts.Count++
			ts.Sum = ts.Sum + t
			ts.Max = utils.If(ts.Max > t, ts.Max, t)
			ts.Min = utils.If(ts.Min < t, ts.Min, t)
			c := utils.ToAny[float64](ts.Count)
			s := utils.ToAny[float64](ts.Sum)
			ts.Average = utils.If(ts.Count > 0, s/c, 0)
		},
		function: func(t *SumStatistics[T]) SumStatistics[T] {
			return *t
		},
	}
}

type Key interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | bool
}

/*
ToMap 切片转Map的Collect方法, 结合stream使用：

eg:

	type TestItemS struct {
		itemNum   int    `json:"itemNum"`
		itemValue string `json:"itemValue"`
	}

res1 := stream.Collect(stream.Of(

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
*/
func ToMap[T any, K Key, U any](keyMapper func(T) K, valueMapper func(T) U, opts ...func(oldV, newV U) U) Collector[T, map[K]U, map[K]U] {
	return &DefaultCollector[T, map[K]U, map[K]U]{
		supplier: func() map[K]U {
			temp := make(map[K]U, 0)
			return temp
		},
		accumulator: func(ts map[K]U, t T) {
			key := keyMapper(t)
			value := valueMapper(t)

			for _, opt := range opts {
				value = mapMerge(key, value, ts, opt)
				ts[key] = value
				return
			}
			ts[key] = value
			return
		},
		function: func(t map[K]U) map[K]U {
			return t
		},
	}
}

func mapMerge[K Key, U any](key K, value U, old map[K]U, opt func(U, U) U) U {
	oldV, ok := old[key]
	var newV U
	if !ok || opt == nil {
		newV = value
	} else {
		newV = opt(oldV, value)
	}
	return newV
}

/*
GroupingBy 分组的Collect 方法，结合stream使用

eg:

res1 := stream.Collect(stream.Of(

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
*/
func GroupingBy[T any, K Key, U any](keyMapper func(T) K, valueMapper func(T) U, opts ...func([]U)) Collector[T, map[K][]U, map[K][]U] {
	return &DefaultCollector[T, map[K][]U, map[K][]U]{
		supplier: func() map[K][]U {
			temp := make(map[K][]U, 0)
			return temp
		},
		accumulator: func(ts map[K][]U, t T) {
			key := keyMapper(t)
			value := valueMapper(t)
			ts[key] = append(ts[key], value)
			return
		},
		function: func(t map[K][]U) map[K][]U {
			for _, vs := range t {
				for _, opt := range opts {
					opt(vs)
				}
			}
			return t
		},
	}
}
