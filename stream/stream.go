// See: zeromicro/go-zero/core/stream
package stream

import (
	"fmt"
	"github.com/todocoder/go-stream/collectors"
	"github.com/todocoder/go-stream/utils"
	"runtime"
	"sort"
	"sync"
)

type (
	Stream[T any] struct {
		source     <-chan T
		isParallel bool
	}

	Optional[T any] struct {
		v *T
	}

	// OptFunc 定义操作切片的函数
	OptFunc[T any] func([]T)
)

func (o Optional[T]) IsPresent() bool {
	return o.v != nil
}

func (o Optional[T]) Get() (v T, ok bool) {
	if o.v == nil {
		return *new(T), false
	}
	return *o.v, true
}
func (o Optional[T]) IfPresent(fn func(T)) {
	if o.v != nil {
		fn(*o.v)
	}
}

func newStream[T any](isParallel bool, items ...T) Stream[T] {
	source := make(chan T, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return Range[T](source, isParallel)
}

func (s Stream[T]) Concat(others ...Stream[T]) Stream[T] {
	source := make(chan T)
	go func() {
		for item := range s.source {
			source <- item
		}
		for _, each := range others {
			each := each
			for item := range each.source {
				source <- item
			}
		}
		close(source)
	}()
	return Range(source, s.isParallel)
}
func (s Stream[T]) Count() (count int64) {
	for range s.source {
		count++
	}
	return
}

func (s Stream[T]) Filter(fn func(item T) bool) Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		if fn(item) {
			pipe <- item
		}
	})
}

func (s Stream[T]) Peek(fn func(item *T)) Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		fn(&item)
		pipe <- item
	})
}

func (s Stream[T]) Limit(maxSize int64) Stream[T] {
	if maxSize < 0 {
		panic("n must not be negative")
	}
	source := make(chan T)
	go func() {
		var n int64 = 0
		for item := range s.source {
			if n < maxSize {
				source <- item
				n++
			} else {
				break
			}
		}
		close(source)
	}()
	return Range(source, s.isParallel)
}

func (s Stream[T]) Skip(n int64) Stream[T] {
	if n < 0 {
		panic("n must not be negative")
	}
	if n == 0 {
		return s
	}
	source := make(chan T)
	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				continue
			} else {
				source <- item
			}
		}
		close(source)
	}()

	return Range(source, s.isParallel)
}

func (s Stream[T]) Distinct(fn func(item T) any) Stream[T] {
	source := make(chan T)
	GoSafe(func() {
		defer close(source)

		keys := make(map[any]struct{})
		for item := range s.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = struct{}{}
			}
		}
	})
	return Range(source, s.isParallel)
}

func (s Stream[T]) Sorted(less func(a, b T) bool) Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
	return Of(items...)
}

func (s Stream[T]) Reverse() Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Of(items...)
}

func (s Stream[T]) Max(comparator func(T, T) int) Optional[T] {
	max, ok := s.FindFirst().Get()
	if !ok {
		return Optional[T]{v: nil}
	}
	s.ForEach(func(t T) {
		if comparator(t, max) > 0 {
			max = t
		}
	})
	return Optional[T]{v: &max}
}

func (s Stream[T]) Min(comparator func(T, T) int) Optional[T] {

	min, ok := s.FindFirst().Get()
	if !ok {
		return Optional[T]{v: nil}
	}
	s.ForEach(func(t T) {
		if comparator(t, min) < 0 {
			min = t
		}
	})
	return Optional[T]{v: &min}
}

type SumIntStatistics[T int | int32 | int64] struct {
	Count   int64
	Sum     int64
	Max     T
	Min     T
	Average float64
}

func (s SumIntStatistics[T]) GetCount() int64 {
	return s.Count
}

func (s SumIntStatistics[T]) GetSum() int64 {
	return s.Sum
}

func (s SumIntStatistics[T]) GetMax() T {
	return s.Max
}

func (s SumIntStatistics[T]) GetMin() T {
	return s.Min
}

func (s SumIntStatistics[T]) GetAverage() float64 {
	return s.Average
}

func (s Stream[T]) SumIntStatistics() SumIntStatistics[int] {
	var cnt = 0
	var sum int64
	var max int
	var min int
	for item := range s.source {
		i := utils.ToAny[int](item)
		if cnt == 0 {
			min = i
		}
		cnt++
		sum = sum + int64(i)
		max = utils.If(max > i, max, i)
		min = utils.If(min < i, min, i)
	}
	return SumIntStatistics[int]{
		Count:   int64(cnt),
		Sum:     sum,
		Max:     max,
		Min:     min,
		Average: utils.If(cnt > 0, float64(sum)/float64(cnt), 0),
	}
}

func (s Stream[T]) SumInt32Statistics() SumIntStatistics[int32] {
	var cnt = 0
	var sum int64
	var max int32
	var min int32
	for item := range s.source {
		i := utils.ToAny[int32](item)
		if cnt == 0 {
			min = i
		}
		cnt++
		sum = sum + int64(i)
		max = utils.If(max > i, max, i)
		min = utils.If(min < i, min, i)
	}
	return SumIntStatistics[int32]{
		Count:   int64(cnt),
		Sum:     sum,
		Max:     max,
		Min:     min,
		Average: utils.If(cnt > 0, float64(sum)/float64(cnt), 0),
	}
}

func (s Stream[T]) SumInt64Statistics() SumIntStatistics[int64] {
	var cnt = 0
	var sum int64
	var max int64
	var min int64
	for item := range s.source {
		i := utils.ToAny[int64](item)
		if cnt == 0 {
			min = i
		}
		cnt++
		sum = sum + i
		max = utils.If(max > i, max, i)
		min = utils.If(min < i, min, i)
	}
	return SumIntStatistics[int64]{
		Count:   int64(cnt),
		Sum:     sum,
		Max:     max,
		Min:     min,
		Average: utils.If(cnt > 0, float64(sum)/float64(cnt), 0),
	}
}

type SumFloatStatistics[T float32 | float64] struct {
	Count   int64
	Sum     float64
	Max     T
	Min     T
	Average float64
}

func (s SumFloatStatistics[T]) GetCount() int64 {
	return s.Count
}
func (s SumFloatStatistics[T]) GetSum() float64 {
	return s.Sum
}
func (s SumFloatStatistics[T]) GetMax() T {
	return s.Max
}
func (s SumFloatStatistics[T]) GetMin() T {
	return s.Min
}
func (s SumFloatStatistics[T]) GetAverage() float64 {
	return s.Average
}

func (s Stream[T]) SumFloat32Statistics() SumFloatStatistics[float32] {
	var cnt = 0
	var sum float64
	var max float32
	var min float32
	for item := range s.source {
		i := utils.ToAny[float32](item)
		if cnt == 0 {
			min = i
		}
		cnt++
		sum = sum + float64(i)
		max = utils.If(max > i, max, i)
		min = utils.If(min < i, min, i)
	}
	return SumFloatStatistics[float32]{
		Count:   int64(cnt),
		Sum:     sum,
		Max:     max,
		Min:     min,
		Average: utils.If(cnt > 0, sum/float64(cnt), 0),
	}
}

func (s Stream[T]) SumFloat64Statistics() SumFloatStatistics[float64] {
	var cnt = 0
	var sum float64
	var max float64
	var min float64
	for item := range s.source {
		i := utils.ToAny[float64](item)
		if cnt == 0 {
			min = i
		}
		cnt++
		sum = sum + i
		max = utils.If(max > i, max, i)
		min = utils.If(min < i, min, i)
	}
	return SumFloatStatistics[float64]{
		Count:   int64(cnt),
		Sum:     sum,
		Max:     max,
		Min:     min,
		Average: utils.If(cnt > 0, sum/float64(cnt), 0),
	}
}

func (s Stream[T]) ForEach(fn func(item T)) {
	var workers = 1
	if s.isParallel {
		workers = runtime.NumCPU() * 2
	}
	//go func() {
	var wg sync.WaitGroup
	// 这里是个占位类型
	pool := make(chan struct{}, workers)
	for item := range s.source {
		val := item
		// 这里是个占位类型值
		pool <- struct{}{}
		wg.Add(1)
		GoSafe(func() {
			defer func() {
				wg.Done()
				<-pool
			}()
			fn(val)
		})
	}
	wg.Wait()
	close(pool)
}

// Walk 让调用者处理每个Item，调用者可以根据给定的Item编写零个、一个或多个项目
func (s Stream[T]) Walk(fn func(item T, pipe chan<- T)) Stream[T] {
	return s.walkLimited(fn)
}

// walkLimited 遍历工作的协程个数限制
func (s Stream[T]) walkLimited(fn func(item T, pipe chan<- T)) Stream[T] {
	var workers = 1
	if s.isParallel {
		workers = runtime.NumCPU() * 2
	}
	pipe := make(chan T, workers)
	go func() {
		defer close(pipe)
		var wg sync.WaitGroup
		// 这里是个占位类型
		pool := make(chan struct{}, workers)
		for item := range s.source {
			val := item
			// 这里是个占位类型值
			pool <- struct{}{}
			wg.Add(1)
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				fn(val, pipe)
			})
		}
		wg.Wait()
	}()
	return Range(pipe, s.isParallel)
}

// AllMatch 返回此流中是否全都满足条件
func (s Stream[T]) AllMatch(predicate func(T) bool) bool {
	// 非缓冲通道
	flag := make(chan bool)
	GoSafe(func() {
		tempFlag := true
		for item := range s.source {
			if !predicate(item) {
				go drain(s.source)
				tempFlag = false
				break
			}
		}
		flag <- tempFlag
	})
	return <-flag
}

// AnyMatch 返回此流中是否存在元素满足所提供的条件
func (s Stream[T]) AnyMatch(predicate func(T) bool) bool {
	flag := make(chan bool)
	GoSafe(func() {
		tempFlag := false
		for item := range s.source {
			if predicate(item) {
				go drain(s.source)
				tempFlag = true
				break
			}
		}
		flag <- tempFlag
	})
	return <-flag
}

// NoneMatch 返回此流中是否全都不满足条件
func (s Stream[T]) NoneMatch(predicate func(T) bool) bool {
	flag := make(chan bool)
	GoSafe(func() {
		tempFlag := true
		for item := range s.source {
			if predicate(item) {
				go drain(s.source)
				tempFlag = false
				break
			}
		}
		flag <- tempFlag
	})

	return <-flag
}

func (s Stream[T]) FindFirst() Optional[T] {
	for item := range s.source {
		go drain(s.source)
		return Optional[T]{v: &item}
	}
	return Optional[T]{v: nil}
}

func (s Stream[T]) FindLast() Optional[T] {
	tempStream := s.Reverse()
	for item := range tempStream.source {
		go drain(tempStream.source)
		return Optional[T]{v: &item}
	}
	return Optional[T]{v: nil}
}

func (s Stream[T]) Reduce(accumulator func(T, T) T) Optional[T] {
	var cnt = 0
	var res T
	for item := range s.source {
		if cnt == 0 {
			cnt++
			res = item
			continue
		}
		cnt++
		res = accumulator(res, item)
	}
	if cnt == 0 {
		return Optional[T]{v: nil}
	}
	return Optional[T]{v: &res}
}

func (s Stream[T]) Map(fn func(item T) any) Stream[any] {
	return Map[T](s, fn)
}

func (s Stream[T]) MapToString(mapper func(T) string) Stream[string] {
	return Map[T](s, mapper)
}

func (s Stream[T]) MapToInt(mapper func(T) int) Stream[int] {
	return Map[T](s, mapper)
}

func (s Stream[T]) MapToInt32(mapper func(T) int32) Stream[int32] {
	return Map[T](s, mapper)
}

func (s Stream[T]) MapToInt64(mapper func(T) int64) Stream[int64] {
	return Map[T](s, mapper)
}

func (s Stream[T]) MapToFloat64(mapper func(T) float64) Stream[float64] {
	return Map[T](s, mapper)
}

func (s Stream[T]) MapToFloat32(mapper func(T) float32) Stream[float32] {
	return Map[T](s, mapper)
}

func (s Stream[T]) FlatMap(mapper func(T) Stream[T]) Stream[T] {
	return FlatMap[T](s, mapper)
}

func (s Stream[T]) FlatMapToString(mapper func(T) Stream[string]) Stream[string] {
	return FlatMap[T](s, mapper)
}

func (s Stream[T]) FlatMapToInt(mapper func(T) Stream[int]) Stream[int] {
	return FlatMap[T](s, mapper)
}
func (s Stream[T]) FlatMapToInt32(mapper func(T) Stream[int32]) Stream[int32] {
	return FlatMap[T](s, mapper)
}
func (s Stream[T]) FlatMapToInt64(mapper func(T) Stream[int64]) Stream[int64] {
	return FlatMap[T](s, mapper)
}
func (s Stream[T]) FlatMapToFloat64(mapper func(T) Stream[float64]) Stream[float64] {
	return FlatMap[T](s, mapper)
}
func (s Stream[T]) FlatMapToFloat32(mapper func(T) Stream[float32]) Stream[float32] {
	return FlatMap[T](s, mapper)
}

func (s Stream[T]) ToSlice() []T {
	r := make([]T, 0)
	for item := range s.source {
		r = append(r, item)
	}
	return r
}

func (s Stream[T]) ToMapString(keyMapper func(T) string, valueMapper func(T) T, opts ...func(oldV, newV T) T) map[string]T {
	res := make(map[string]T, 0)
	for item := range s.source {
		key := keyMapper(item)
		value := valueMapper(item)

		oldV, ok := res[key]

		if !ok || opts == nil || len(opts) < 1 {
			res[key] = value
			continue
		}
		var newV T
		for _, opt := range opts {
			if !ok || opt == nil {
				newV = value
			} else {
				newV = opt(oldV, value)
			}
			res[key] = newV
		}
	}
	return res
}

func (s Stream[T]) ToMapInt(keyMapper func(T) int, valueMapper func(T) T, opts ...func(oldV, newV T) T) map[int]T {
	res := make(map[int]T, 0)
	for item := range s.source {
		key := keyMapper(item)
		value := valueMapper(item)

		oldV, ok := res[key]
		if !ok || opts == nil || len(opts) < 1 {
			res[key] = value
			continue
		}
		var newV T
		for _, opt := range opts {
			if !ok || opt == nil {
				newV = value
			} else {
				newV = opt(oldV, value)
			}
			res[key] = newV
		}
	}
	return res
}

// Deprecated: This function is no longer recommended. Please use extracted.Collect() instead.
func (s Stream[T]) Collect(collector collectors.Collector[T, T, any]) any {
	temp := collector.Supplier()()
	for item := range s.source {
		collector.Accumulator()(item, temp)
	}
	return collector.Finisher()(temp)
}

func drain[T any](channel <-chan T) {
	for range channel {
	}
}

func Of[T any](values ...T) Stream[T] {
	return newStream(false, values...)
}

func OfParallel[T any](values ...T) Stream[T] {
	s := newStream(true, values...)
	return s.Walk(func(item T, pipe chan<- T) {
		pipe <- item
	})
}
func OfFrom[T any](generate func(source chan<- T)) Stream[T] {
	source := make(chan T)
	GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range[T](source, false)
}

func OfFromParallel[T any](generate func(source chan<- T)) Stream[T] {
	source := make(chan T)
	GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range[T](source, true)
}

func Range[T any](source <-chan T, isParallel bool) Stream[T] {
	return Stream[T]{
		source:     source,
		isParallel: isParallel,
	}
}

func (s Stream[T]) GroupingByString(groupFunc func(T) string, opts ...OptFunc[T]) map[string][]T {
	groups := make(map[string][]T)
	s.ForEach(func(t T) {
		key := groupFunc(t)
		groups[key] = append(groups[key], t)
	})
	for _, vs := range groups {
		for _, opt := range opts {
			opt(vs)
		}
	}
	return groups
}

func (s Stream[T]) GroupingByInt(groupFunc func(T) int, opts ...OptFunc[T]) map[int][]T {
	groups := make(map[int][]T)
	s.ForEach(func(t T) {
		key := groupFunc(t)
		groups[key] = append(groups[key], t)
	})
	for _, vs := range groups {
		for _, opt := range opts {
			opt(vs)
		}
	}
	return groups
}

// Concat 拼接流
func Concat[T any](s Stream[T], others ...Stream[T]) Stream[T] {
	return s.Concat(others...)
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		fmt.Println(p)
	}
}
