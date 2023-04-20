/*
 * Copyright 2023 zerune.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package function

// Runnable
// See: java.lang.Runnable
type Runnable func()

// Comparator
// See: java.util.Comparator
type Comparator[T any] func(T, T) int

// BiConsumer
// See: java.util.function.BiConsumer
type BiConsumer[T any, U any] func(T, U)

// BiFunction
// See: java.util.function.BiFunction
type BiFunction[T any, U any, R any] func(T, U) R

// BinaryOperator
// See: java.util.function.BinaryOperator
type BinaryOperator[T any] BiFunction[T, T, T]

// BiPredicate
// See: java.util.function.BiPredicate
type BiPredicate[T any, U any] func(T, U) bool

// BooleanSupplier
// See: java.util.function.BooleanSupplier
type BooleanSupplier func() bool

// Consumer
// See: java.util.function.Consumer
type Consumer[T any] func(T)

// DoubleBinaryOperator
// See: java.util.function.DoubleBinaryOperator
type DoubleBinaryOperator func(float64, float64) float64

// DoubleConsumer
// See: java.util.function.DoubleConsumer
type DoubleConsumer func(float64)

// DoubleFunction
// See: java.util.function.DoubleFunction
type DoubleFunction[R any] func(float64) R

// DoublePredicate
// See: java.util.function.DoublePredicate
type DoublePredicate func(float64) bool

// DoubleSupplier
// See: java.util.function.DoubleSupplier
type DoubleSupplier func() float64

// DoubleToIntFunction
// See: java.util.function.DoubleToIntFunction
type DoubleToIntFunction func(float64) int

// DoubleToLongFunction
// See: java.util.function.DoubleToLongFunction
type DoubleToLongFunction func(float64) int64

// DoubleUnaryOperator
// See: java.util.function.DoubleUnaryOperator
type DoubleUnaryOperator func(float64) float64

// Function
// See: java.util.function.Function
type Function[T any, R any] func(T) R

// IntBinaryOperator
// See: java.util.function.IntBinaryOperator
type IntBinaryOperator func(int, int) int

// IntConsumer
// See: java.util.function.IntConsumer
type IntConsumer func(int)

// IntFunction
// See: java.util.function.IntFunction
type IntFunction[R any] func(int) R

// IntPredicate
// See: java.util.function.IntPredicate
type IntPredicate func(int) bool

// IntSupplier
// See: java.util.function.IntSupplier
type IntSupplier func() int

// IntToDoubleFunction
// See: java.util.function.IntToDoubleFunction
type IntToDoubleFunction func(int) float64

// IntToLongFunction
// See: java.util.function.IntToLongFunction
type IntToLongFunction func(int) int64

// IntUnaryOperator
// See: java.util.function.IntUnaryOperator
type IntUnaryOperator func(int) int

// LongBinaryOperator
// See: java.util.function.LongBinaryOperator
type LongBinaryOperator func(int64, int64) int64

// LongConsumer
// See: java.util.function.LongConsumer
type LongConsumer func(int64)

// LongFunction
// See: java.util.function.LongFunction
type LongFunction[R any] func(int64) R

// LongPredicate
// See: java.util.function.LongPredicate
type LongPredicate func(int64) bool

// LongSupplier
// See: java.util.function.LongSupplier
type LongSupplier func() int64

// LongToDoubleFunction
// See: java.util.function.LongToDoubleFunction
type LongToDoubleFunction func(int64) float64

// LongToIntFunction
// See: java.util.function.LongToIntFunction
type LongToIntFunction func(int64) int

// LongUnaryOperator
// See: java.util.function.LongUnaryOperator
type LongUnaryOperator func(int64) int64

// ObjDoubleConsumer
// See: java.util.function.ObjDoubleConsumer
type ObjDoubleConsumer[T any] func(T, float64)

// ObjIntConsumer
// See: java.util.function.ObjIntConsumer
type ObjIntConsumer[T any] func(T, int)

// ObjLongConsumer
// See: java.util.function.ObjLongConsumer
type ObjLongConsumer[T any] func(T, int64)

// Predicate
// See: java.util.function.Predicate
type Predicate[T any] func(T) bool

// Supplier
// See: java.util.function.Supplier
type Supplier[T any] func() T

// ToDoubleFunction
// See: java.util.function.ToDoubleFunction
type ToDoubleFunction[T any] func(T) float64

// ToDoubleBiFunction
// See: java.util.function.ToDoubleBiFunction
type ToDoubleBiFunction[T any, U any] func(T, U) float64

// ToIntFunction
// See: java.util.function.ToIntFunction
type ToIntFunction[T any] func(T) int

// ToIntBiFunction
// See: java.util.function.ToIntBiFunction
type ToIntBiFunction[T any, U any] func(T, U) int

// ToLongFunction
// See: java.util.function.ToLongFunction
type ToLongFunction[T any] func(T) int64

// ToLongBiFunction
// See: java.util.function.ToLongBiFunction
type ToLongBiFunction[T any, U any] func(T, U) int64

// UnaryOperator
// See: java.util.function.UnaryOperator
type UnaryOperator[T any] Function[T, T]
