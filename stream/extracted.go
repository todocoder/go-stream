package stream

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
