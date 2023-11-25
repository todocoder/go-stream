package stream

import "github.com/todocoder/go-stream/collectors"

/*
Map stream 流 类型转换方法

eg:

res := Map(Of(

		TestItem{itemNum: 1, itemValue: "item1"},
		TestItem{itemNum: 2, itemValue: "item2"},
		TestItem{itemNum: 3, itemValue: "item3"},
	), func(item TestItem) ToTestItem {
		return ToTestItem{
			itemNum:   item.itemNum,
			itemValue: item.itemValue,
		}
	}).ToSlice()
	fmt.Println(res)
*/
func Map[T any, R any](s Stream[T], mapper func(T) R) Stream[R] {
	mapped := make([]R, 0)

	for el := range s.source {
		mapped = append(mapped, mapper(el))
	}
	return Of(mapped...)
}

func FlatMap[T any, R any](s Stream[T], mapper func(T) Stream[R]) Stream[R] {
	streams := make([]Stream[R], 0)
	s.ForEach(func(t T) {
		streams = append(streams, mapper(t))
	})

	newEl := make([]R, 0)
	for _, str := range streams {
		newEl = append(newEl, str.ToSlice()...)
	}

	return Of(newEl...)
}

func GroupingBy[T any, K string | int | int32 | int64, R any](s Stream[T], keyMapper func(T) K, valueMapper func(T) R, opts ...OptFunc[R]) map[K][]R {
	groups := make(map[K][]R)
	s.ForEach(func(t T) {
		key := keyMapper(t)
		groups[key] = append(groups[key], valueMapper(t))
	})
	for _, vs := range groups {
		for _, opt := range opts {
			opt(vs)
		}
	}
	return groups
}

func Collect[T any, A any, R any](s Stream[T], collector collectors.Collector[T, A, R]) R {
	temp := collector.Supplier()()
	for item := range s.source {
		collector.Accumulator()(temp, item)
	}
	return collector.Finisher()(temp)
}
