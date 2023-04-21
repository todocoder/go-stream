package stream

import (
	"runtime"
	"sort"

	"github.com/zerune/go-core/lang"
	"github.com/zerune/go-core/routine"
	"github.com/zerune/go-core/routine/channel"
	"github.com/zerune/go-core/util/function"
	"github.com/zerune/go-core/util/optional"
	"github.com/zerune/go-core/util/stream/collectors"
)

type Stream[T any] interface {
	Concat(others ...Stream[T]) Stream[T]
	Count() int64
	Filter(predicate function.Predicate[T]) Stream[T]
	Peek(action function.Consumer[*T]) Stream[T]
	Limit(maxSize int64) Stream[T]
	Skip(n int64) Stream[T]
	Distinct(opts ...function.Function[T, any]) Stream[T]
	Sorted(comparator ...function.Comparator[T]) Stream[T]
	Reverse() Stream[T]
	Max(comparator function.Comparator[T]) optional.Optional[T]
	Min(comparator function.Comparator[T]) optional.Optional[T]
	ForEach(action function.Consumer[T])
	AllMatch(predicate function.Predicate[T]) bool
	AnyMatch(predicate function.Predicate[T]) bool
	NoneMatch(predicate function.Predicate[T]) bool
	TakeWhile(predicate function.Predicate[T]) Stream[T]
	DropWhile(predicate function.Predicate[T]) Stream[T]
	FindAny() optional.Optional[T]
	FindFirst() optional.Optional[T]
	FindLast() optional.Optional[T]
	Reduce(accumulator function.BinaryOperator[T]) optional.Optional[T]
	Map(mapper function.Function[T, any]) Stream[any]
	MapToInt(mapper function.ToIntFunction[T]) Stream[int]
	MapToLong(mapper function.ToLongFunction[T]) Stream[int64]
	MapToDouble(mapper function.ToDoubleFunction[T]) Stream[float64]
	FlatMap(mapper function.Function[T, Stream[any]]) Stream[any]
	FlatMapToInt(mapper function.Function[T, Stream[int]]) Stream[int]
	FlatMapToLong(mapper function.Function[T, Stream[int64]]) Stream[int64]
	FlatMapToDouble(mapper function.Function[T, Stream[float64]]) Stream[float64]
	ToSlice() []T
	Collect(collector collectors.Collector[T, any, any]) any
}

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers          int
	}

	// walkFunc defines the method to walk through all the elements in a Stream.
	walkFunc[T any] func(item T, pipe chan<- T)

	streamImpl[T any] struct {
		source     <-chan T
		isParallel bool
	}
)

var (
	defaultWorkers = 1
	maxWorkers     = runtime.NumCPU() * 2
)

func (s *streamImpl[T]) Concat(others ...Stream[T]) Stream[T] {
	source := make(chan T)
	go func() {
		if s.isParallel {
			group := routine.NewRoutineGroup()
			group.Run(
				func() {
					for item := range s.source {
						source <- item
					}
				},
			)
			group.Run(
				func() {
					for _, each := range others {
						each.ForEach(func(item T) {
							source <- item
						})
					}
				},
			)
			group.Wait()
		} else {
			for item := range s.source {
				source <- item
			}
			for _, each := range others {
				each.ForEach(func(item T) {
					source <- item
				})
			}
		}
		close(source)
	}()
	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) Count() (count int64) {
	for range s.source {
		count++
	}
	return
}

func (s *streamImpl[T]) Filter(predicate function.Predicate[T]) Stream[T] {
	return s.walk(func(item T, pipe chan<- T) {
		if predicate(item) {
			pipe <- item
		}
	})
}

func (s *streamImpl[T]) Peek(action function.Consumer[*T]) Stream[T] {
	return s.walk(func(item T, pipe chan<- T) {
		action(&item)
		pipe <- item
	})
}

func (s *streamImpl[T]) Limit(maxSize int64) Stream[T] {
	if maxSize < 0 {
		maxSize = 0
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
	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) Skip(n int64) Stream[T] {
	if n <= 0 {
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

	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) Distinct(opts ...function.Function[T, any]) Stream[T] {
	source := make(chan T)
	routine.GoSafe(func() {
		defer close(source)

		keys := make(map[any]lang.PlaceholderType)
		for item := range s.source {
			var key any
			key = item
			for _, opt := range opts {
				key = opt(item)
			}
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = lang.PlaceholderType{}
			}
		}
	})
	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) Sorted(comparator ...function.Comparator[T]) Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}

	flag := true
	for _, item := range items {
		var v any
		v = &item
		if _, ok := v.(lang.Comparable[T]); !ok {
			flag = false
			break
		}
	}

	if flag {
		sort.Slice(items, func(i, j int) bool {
			var a any
			a = &items[i]
			v, _ := a.(lang.Comparable[T])
			return v.CompareTo(items[j]) > 0
		})
	} else {
		for _, comp := range comparator {
			sort.Slice(items, func(i, j int) bool {
				return comp(items[i], items[j]) > 0
			})
		}
	}

	return Of(items...)
}

func (s *streamImpl[T]) Reverse() Stream[T] {
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

func (s *streamImpl[T]) Max(comparator function.Comparator[T]) optional.Optional[T] {
	first := s.FindFirst()
	if !first.IsPresent() {
		return optional.Empty[T]()
	}
	max := *first.Get()
	s.ForEach(func(t T) {
		if comparator(t, max) > 0 {
			max = t
		}
	})
	return optional.Of(max)
}

func (s *streamImpl[T]) Min(comparator function.Comparator[T]) optional.Optional[T] {
	first := s.FindFirst()
	if !first.IsPresent() {
		return optional.Empty[T]()
	}
	min := *first.Get()
	s.ForEach(func(t T) {
		if comparator(t, min) < 0 {
			min = t
		}
	})
	return optional.Of(min)
}

func (s *streamImpl[T]) ForEach(action function.Consumer[T]) {
	option := buildOptions(s.isParallel)
	wg := routine.NewLimitedGroup(option.workers)
	for item := range s.source {
		// important, used in another goroutine
		val := item

		// better to safely run caller defined method
		wg.RunSafe(func() {
			action(val)
		})
	}
	wg.Wait()
}

// AllMatch 返回此流中是否全都满足条件
func (s *streamImpl[T]) AllMatch(predicate function.Predicate[T]) bool {
	// 非缓冲通道
	flag := make(chan bool)
	routine.GoSafe(func() {
		tempFlag := true
		for item := range s.source {
			if !predicate(item) {
				go channel.Drain(s.source)
				tempFlag = false
				break
			}
		}
		flag <- tempFlag
	})
	return <-flag
}

// AnyMatch 返回此流中是否存在元素满足所提供的条件
func (s *streamImpl[T]) AnyMatch(predicate function.Predicate[T]) bool {
	flag := make(chan bool)
	routine.GoSafe(func() {
		tempFlag := false
		for item := range s.source {
			if predicate(item) {
				go channel.Drain(s.source)
				tempFlag = true
				break
			}
		}
		flag <- tempFlag
	})

	return <-flag
}

// NoneMatch 返回此流中是否全都不满足条件
func (s *streamImpl[T]) NoneMatch(predicate function.Predicate[T]) bool {
	flag := make(chan bool)
	routine.GoSafe(func() {
		tempFlag := true
		for item := range s.source {
			if predicate(item) {
				go channel.Drain(s.source)
				tempFlag = false
				break
			}
		}
		flag <- tempFlag
	})

	return <-flag
}

func (s *streamImpl[T]) TakeWhile(predicate function.Predicate[T]) Stream[T] {
	source := make(chan T)
	go func() {
		for item := range s.source {
			if predicate(item) {
				source <- item
			} else {
				go channel.Drain(s.source)
				break
			}
		}
		close(source)
	}()

	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) DropWhile(predicate function.Predicate[T]) Stream[T] {
	source := make(chan T)
	go func() {
		drop := true
		for item := range s.source {
			if predicate(item) && drop {
				continue
			} else {
				if drop {
					drop = false
				}
				source <- item
			}
		}
		close(source)
	}()

	return newRangeStream(source, s.isParallel)
}

func (s *streamImpl[T]) FindAny() optional.Optional[T] {
	return s.FindFirst()
}

func (s *streamImpl[T]) FindFirst() optional.Optional[T] {
	for item := range s.source {
		go channel.Drain(s.source)
		return optional.Of(item)
	}
	return optional.Empty[T]()
}

func (s *streamImpl[T]) FindLast() optional.Optional[T] {
	tempStream := s.Reverse().(*streamImpl[T])
	for item := range tempStream.source {
		go channel.Drain(tempStream.source)
		return optional.Of(item)
	}
	return optional.Empty[T]()
}

func (s *streamImpl[T]) Reduce(accumulator function.BinaryOperator[T]) optional.Optional[T] {
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
		return optional.Empty[T]()
	}
	return optional.Of(res)
}

func (s *streamImpl[T]) Map(mapper function.Function[T, any]) Stream[any] {
	return Map[T](s, mapper)
}

func (s *streamImpl[T]) MapToInt(mapper function.ToIntFunction[T]) Stream[int] {
	return Map[T](s, mapper)
}

func (s *streamImpl[T]) MapToLong(mapper function.ToLongFunction[T]) Stream[int64] {
	return Map[T](s, mapper)
}

func (s *streamImpl[T]) MapToDouble(mapper function.ToDoubleFunction[T]) Stream[float64] {
	return Map[T](s, mapper)
}

func (s *streamImpl[T]) FlatMap(mapper function.Function[T, Stream[any]]) Stream[any] {
	return FlatMap[T](s, mapper)
}

func (s *streamImpl[T]) FlatMapToInt(mapper function.Function[T, Stream[int]]) Stream[int] {
	return FlatMap[T](s, mapper)
}

func (s *streamImpl[T]) FlatMapToLong(mapper function.Function[T, Stream[int64]]) Stream[int64] {
	return FlatMap[T](s, mapper)
}

func (s *streamImpl[T]) FlatMapToDouble(mapper function.Function[T, Stream[float64]]) Stream[float64] {
	return FlatMap[T](s, mapper)
}

func (s *streamImpl[T]) ToSlice() []T {
	r := make([]T, 0)
	for item := range s.source {
		r = append(r, item)
	}
	return r
}

func (s *streamImpl[T]) Collect(collector collectors.Collector[T, any, any]) any {
	container := collector.Supplier()()
	for item := range s.source {
		container = collector.Accumulator()(container, item)
	}
	return collector.Finisher()(container)
}

// lets the callers handle each item, the caller may write zero, one or more items base on the given item.
func (s *streamImpl[T]) walk(fn walkFunc[T]) Stream[T] {
	option := buildOptions(s.isParallel)
	if option.unlimitedWorkers {
		wg := routine.NewRoutineGroup()
		return s.walkFunc(wg, fn, option)
	}
	wg := routine.NewLimitedGroup(option.workers)
	return s.walkFunc(wg, fn, option)
}

func (s *streamImpl[T]) walkFunc(wg routine.WaitGroup, fn walkFunc[T], option rxOptions) Stream[T] {
	pipe := make(chan T, option.workers)

	go func() {
		for item := range s.source {
			// important, used in another goroutine
			val := item
			// better to safely run caller defined method
			wg.RunSafe(
				func() {
					fn(val, pipe)
				},
			)
		}

		wg.Wait()
		close(pipe)
	}()

	return newRangeStream(pipe, s.isParallel)
}

func Of[T any](values ...T) Stream[T] {
	return newStream(false, values...)
}

func OfParallel[T any](values ...T) Stream[T] {
	return newStream(true, values...)
}

func OfFrom[T any](generate func(source chan<- T)) Stream[T] {
	source := make(chan T)
	routine.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return newRangeStream[T](source, false)
}

func OfFromParallel[T any](generate func(source chan<- T)) Stream[T] {
	source := make(chan T)
	routine.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return newRangeStream[T](source, true)
}

// Concat returns a concatenated Stream.
func Concat[T any](s Stream[T], others ...Stream[T]) Stream[T] {
	return s.Concat(others...)
}

func newStream[T any](isParallel bool, items ...T) Stream[T] {
	source := make(chan T, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return newRangeStream[T](source, isParallel)
}

func newRangeStream[T any](source <-chan T, isParallel bool) Stream[T] {
	return &streamImpl[T]{
		source:     source,
		isParallel: isParallel,
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(parallel bool) rxOptions {
	options := rxOptions{
		workers: defaultWorkers,
	}
	if parallel {
		options.workers = maxWorkers
	}
	return options
}
