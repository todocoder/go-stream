/**
* See java: stream
 */
package collectors

type Supplier[T any] func() T

type Function[T any, R any] func(T) R

type BiConsumer[T any, R any] func(T, R) R

type BinaryOperator[T any, A any, R any] func(T, A) R

type Collector[T any, A any, R any] interface {
	Supplier() Supplier[R]
	Accumulator() BiConsumer[T, R]
	Combiner() BiConsumer[R, R]
	Finisher() Function[R, any]
	BinaryOperator() BinaryOperator[T, A, R]
}

type DefaultCollector[T any, A any, R any] struct {
	supplier       Supplier[R]
	accumulator    BiConsumer[T, R]
	combiner       BiConsumer[R, R]
	function       Function[R, R]
	binaryOperator BinaryOperator[T, A, R]
}

func (s *DefaultCollector[T, A, R]) Supplier() Supplier[R]                   { return s.supplier }
func (s *DefaultCollector[T, A, R]) Accumulator() BiConsumer[T, R]           { return s.accumulator }
func (s *DefaultCollector[T, A, R]) Combiner() BiConsumer[R, R]              { return s.combiner }
func (s *DefaultCollector[T, A, R]) Finisher() Function[R, R]                { return s.function }
func (s *DefaultCollector[T, A, R]) BinaryOperator() BinaryOperator[T, A, R] { return s.binaryOperator }

func MapMerge(key any, value any, old map[any]any, opt func(any, any) any) any {
	oldV := old[key]
	var newV any
	if oldV == nil || opt == nil {
		newV = value
	} else {
		newV = opt(oldV, value)
	}
	return newV
}

func MapMergeS(key string, value any, old map[string]any, opt func(any, any) any) any {
	oldV := old[key]
	var newV any
	if oldV == nil || opt == nil {
		newV = value
	} else {
		newV = opt(oldV, value)
	}
	return newV
}
