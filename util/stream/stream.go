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

var (
	defaultWorkers = 1
	maxWorkers     = runtime.NumCPU() * 2
)

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers          int
	}

	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc[T any] func(item T, pipe chan<- T)
	Stream[T any]   struct {
		source     <-chan T
		isParallel bool
	}
)

func (s *Stream[T]) Concat(others ...*Stream[T]) *Stream[T] {
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
						for item := range each.source {
							source <- item
						}
					}
				},
			)
			group.Wait()
		} else {
			for item := range s.source {
				source <- item
			}
			for _, each := range others {
				for item := range each.source {
					source <- item
				}
			}
		}
		close(source)
	}()
	return newRangeStream(source, s.isParallel)
}

func (s *Stream[T]) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

func (s *Stream[T]) Filter(predicate function.Predicate[T]) *Stream[T] {
	return s.walk(func(item T, pipe chan<- T) {
		if predicate(item) {
			pipe <- item
		}
	})
}

func (s *Stream[T]) Peek(action function.Consumer[*T]) *Stream[T] {
	return s.walk(func(item T, pipe chan<- T) {
		action(&item)
		pipe <- item
	})
}

func (s *Stream[T]) Limit(maxSize int64) *Stream[T] {
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

func (s *Stream[T]) Skip(n int64) *Stream[T] {
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

func (s *Stream[T]) Distinct(opts ...function.Function[T, any]) *Stream[T] {
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

func (s *Stream[T]) Sorted(comparator ...function.Comparator[T]) *Stream[T] {
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

func (s *Stream[T]) Reverse() *Stream[T] {
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

func (s *Stream[T]) Max(comparator func(T, T) int) *optional.Optional[T] {
	max, ok := s.FindFirst().Get()
	if !ok {
		return optional.OfNullable[T](nil)
	}
	s.ForEach(func(t T) {
		if comparator(t, max) > 0 {
			max = t
		}
	})
	return optional.Of(max)
}

func (s *Stream[T]) Min(comparator func(T, T) int) *optional.Optional[T] {
	min, ok := s.FindFirst().Get()
	if !ok {
		return optional.OfNullable[T](nil)
	}
	s.ForEach(func(t T) {
		if comparator(t, min) < 0 {
			min = t
		}
	})
	return optional.Of(min)
}

func (s *Stream[T]) ForEach(action function.Consumer[T]) {
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
func (s *Stream[T]) AllMatch(predicate function.Predicate[T]) bool {
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
func (s *Stream[T]) AnyMatch(predicate function.Predicate[T]) bool {
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
func (s *Stream[T]) NoneMatch(predicate function.Predicate[T]) bool {
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

func (s *Stream[T]) TakeWhile(predicate function.Predicate[T]) *Stream[T] {
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

func (s *Stream[T]) DropWhile(predicate function.Predicate[T]) *Stream[T] {
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

func (s *Stream[T]) FindAny() *optional.Optional[T] {
	return s.FindFirst()
}

func (s *Stream[T]) FindFirst() *optional.Optional[T] {
	for item := range s.source {
		go channel.Drain(s.source)
		return optional.Of(item)
	}
	return optional.OfNullable[T](nil)
}

func (s *Stream[T]) FindLast() *optional.Optional[T] {
	tempStream := s.Reverse()
	for item := range tempStream.source {
		go channel.Drain(tempStream.source)
		return optional.Of(item)
	}
	return optional.OfNullable[T](nil)
}

func (s *Stream[T]) Reduce(accumulator function.BinaryOperator[T]) *optional.Optional[T] {
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
		return optional.OfNullable[T](nil)
	}
	return optional.Of(res)
}

func (s *Stream[T]) Map(fn func(item T) any) *Stream[any] {
	return Map[T](s, fn)
}

func (s *Stream[T]) MapToInt(mapper func(T) int64) *Stream[int64] {
	return Map[T](s, mapper)
}

func (s *Stream[T]) MapToDouble(mapper func(T) float64) *Stream[float64] {
	return Map[T](s, mapper)
}

func (s *Stream[T]) FlatMap(mapper func(T) *Stream[any]) *Stream[any] {
	return FlatMap[T](s, mapper)
}

func (s *Stream[T]) FlatMapToInt(mapper func(T) *Stream[any]) *Stream[any] {
	return FlatMap[T](s, mapper)
}

func (s *Stream[T]) FlatMapToDouble(mapper func(T) *Stream[float64]) *Stream[float64] {
	return FlatMap[T](s, mapper)
}

func (s *Stream[T]) ToSlice() []T {
	r := make([]T, 0)
	for item := range s.source {
		r = append(r, item)
	}
	return r
}

func (s *Stream[T]) Collect(collector collectors.Collector[T, any, any]) any {
	container := collector.Supplier()()
	for item := range s.source {
		container = collector.Accumulator()(container, item)
	}
	return collector.Finisher()(container)
}

// lets the callers handle each item, the caller may write zero, one or more items base on the given item.
func (s *Stream[T]) walk(fn WalkFunc[T]) *Stream[T] {
	option := buildOptions(s.isParallel)
	if option.unlimitedWorkers {
		wg := routine.NewRoutineGroup()
		return s.walkFunc(wg, fn, option)
	}
	wg := routine.NewLimitedGroup(option.workers)
	return s.walkFunc(wg, fn, option)
}

func (s *Stream[T]) walkFunc(wg routine.WaitGroup, fn WalkFunc[T], option *rxOptions) *Stream[T] {
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

func Of[T any](values ...T) *Stream[T] {
	return newStream(false, values...)
}

func OfParallel[T any](values ...T) *Stream[T] {
	return newStream(true, values...)
}

func OfFrom[T any](generate func(source chan<- T)) *Stream[T] {
	source := make(chan T)
	routine.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return newRangeStream[T](source, false)
}

func OfFromParallel[T any](generate func(source chan<- T)) *Stream[T] {
	source := make(chan T)
	routine.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return newRangeStream[T](source, true)
}

// Concat returns a concatenated Stream.
func Concat[T any](s *Stream[T], others ...*Stream[T]) *Stream[T] {
	return s.Concat(others...)
}

func newStream[T any](isParallel bool, items ...T) *Stream[T] {
	source := make(chan T, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return newRangeStream[T](source, isParallel)
}

func newRangeStream[T any](source <-chan T, isParallel bool) *Stream[T] {
	return &Stream[T]{
		source:     source,
		isParallel: isParallel,
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(parallel bool) *rxOptions {
	options := &rxOptions{
		workers: defaultWorkers,
	}
	if parallel {
		options.workers = maxWorkers
	}
	return options
}
