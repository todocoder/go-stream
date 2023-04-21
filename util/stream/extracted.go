package stream

func Map[T any, R any](s Stream[T], mapper func(T) R) Stream[R] {
	mapped := make([]R, 0)
	s.ForEach(func(item T) {
		mapped = append(mapped, mapper(item))
	})
	return Of(mapped...)
}

func FlatMap[T any, R any](s Stream[T], mapper func(T) Stream[R]) Stream[R] {
	streams := make([]Stream[R], 0)
	s.ForEach(func(t T) {
		streams = append(streams, mapper(t))
	})

	newEl := make([]R, 0)
	for _, item := range streams {
		newEl = append(newEl, item.ToSlice()...)
	}

	return Of(newEl...)
}
