package collectors

import "github.com/zerune/go-core/util/function"

func ToSlice[T any]() Collector[T, any, any] {
	return &DefaultCollector[T, any, any]{
		supplier: func() any {
			return make([]T, 0)
		},
		accumulator: func(container any, u T) any {
			container = append(container.([]T), u)
			return container
		},
		function: func(t any) any {
			return t
		},
	}
}

func ToMap[T any](
	keyMapper function.Function[T, any], valueMapper function.Function[T, any],
	mergeFunction ...function.BinaryOperator[any],
) Collector[T, any, any] {
	return &DefaultCollector[T, any, any]{
		supplier: func() any {
			return make(map[any]any, 0)
		},
		accumulator: func(container any, u T) any {
			key := keyMapper(u)
			value := valueMapper(u)
			c := container.(map[any]any)
			for _, opt := range mergeFunction {
				value = mapMerge(key, value, c, opt)
				c[key] = value
				return c
			}
			c[key] = value
			return c
		},
		function: func(t any) any {
			return t
		},
	}
}

func GroupingBy[T any](keyMapper function.Function[T, any], valueMapper function.Function[T, any]) Collector[T, any, any] {
	return &DefaultCollector[T, any, any]{
		supplier: func() any {
			temp := make(map[any][]any)
			return temp
		},
		accumulator: func(container any, u T) any {
			key := keyMapper(u)
			value := valueMapper(u)
			c := container.(map[any][]any)
			c[key] = append(c[key], value)
			return c
		},
		function: func(t any) any {
			return t
		},
	}
}

func mapMerge(key any, value any, container map[any]any, mergeFunction function.BinaryOperator[any]) any {
	oldV := container[key]
	var newV any
	if oldV == nil || mergeFunction == nil {
		newV = value
	} else {
		newV = mergeFunction(oldV, value)
	}
	return newV
}
