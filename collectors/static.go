/**
* See java: stream
 */
package collectors

type Supplier[R any] func() R

type Function[A any, R any] func(A) R

type BiConsumer[A any, T any] func(A, T)

type BinaryOperator[T any, A any, R any] func(T, A) R

// Collector See java: stream
// <T> – the type of input elements to the reduction operation
// <A> – the mutable accumulation type of the reduction operation (often hidden as an implementation detail)
// <R> – the result type of the reduction operation
type Collector[T any, A any, R any] interface {
	Supplier() Supplier[A]
	Accumulator() BiConsumer[A, T]
	Combiner() BiConsumer[R, R]
	Finisher() Function[A, R]
	BinaryOperator() BinaryOperator[T, A, R]
}

type DefaultCollector[T any, A any, R any] struct {
	supplier       Supplier[A]
	accumulator    BiConsumer[A, T]
	combiner       BiConsumer[R, R]
	function       Function[A, R]
	binaryOperator BinaryOperator[T, A, R]
}

func (s *DefaultCollector[T, A, R]) Supplier() Supplier[A]                   { return s.supplier }
func (s *DefaultCollector[T, A, R]) Accumulator() BiConsumer[A, T]           { return s.accumulator }
func (s *DefaultCollector[T, A, R]) Combiner() BiConsumer[R, R]              { return s.combiner }
func (s *DefaultCollector[T, A, R]) Finisher() Function[A, R]                { return s.function }
func (s *DefaultCollector[T, A, R]) BinaryOperator() BinaryOperator[T, A, R] { return s.binaryOperator }
