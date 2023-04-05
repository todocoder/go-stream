package collectors

func ToSlice[T any]() Collector[T, T, any] {
	return &DefaultCollector[T, T, any]{
		supplier: func() any {
			temp := make([]T, 0)
			return temp
		},
		accumulator: func(t T, ts any) any {
			ts = append(ts.([]T), t)
			return ts
		},
		function: func(t any) any {
			return t
		},
	}
}
func ToMap[T any](keyMapper func(T) any, valueMapper func(T) any, opts ...func(oldV, newV any) any) Collector[T, T, any] {
	return &DefaultCollector[T, T, any]{
		supplier: func() any {
			temp := make(map[any]any, 0)
			return temp
		},
		accumulator: func(t T, ts any) any {
			key := keyMapper(t)
			value := valueMapper(t)

			for _, opt := range opts {
				value = MapMerge(key, value, ts.(map[any]any), opt)
				(ts.(map[any]any))[key] = value
				return ts
			}
			(ts.(map[any]any))[key] = value
			return ts
		},
		function: func(t any) any {
			return t
		},
	}
}

func GroupingBy[T any](keyMapper func(T) any, valueMapper func(T) any) Collector[T, T, any] {
	return &DefaultCollector[T, T, any]{
		supplier: func() any {
			temp := make(map[any][]any)
			return temp
		},
		accumulator: func(t T, ts any) any {
			key := keyMapper(t)
			value := valueMapper(t)
			(ts.(map[any][]any))[key] = append((ts.(map[any][]any))[key], value)
			return ts
		},
		function: func(t any) any {
			return t
		},
	}
}
