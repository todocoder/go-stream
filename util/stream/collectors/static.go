package collectors

import "github.com/zerune/go-core/util/function"

type Collector[T any, A any, R any] interface {
	Supplier() function.Supplier[A]
	Accumulator() function.BiFunction[A, T, A]
	Combiner() function.BinaryOperator[A]
	Finisher() function.Function[A, R]
}

type DefaultCollector[T any, A any, R any] struct {
	supplier    function.Supplier[A]
	accumulator function.BiFunction[A, T, A]
	combiner    function.BinaryOperator[A]
	function    function.Function[A, R]
}

func (s *DefaultCollector[T, A, R]) Supplier() function.Supplier[A]            { return s.supplier }
func (s *DefaultCollector[T, A, R]) Accumulator() function.BiFunction[A, T, A] { return s.accumulator }
func (s *DefaultCollector[T, A, R]) Combiner() function.BinaryOperator[A]      { return s.combiner }
func (s *DefaultCollector[T, A, R]) Finisher() function.Function[A, R]         { return s.function }
