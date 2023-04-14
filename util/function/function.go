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

// Function
// See: java.util.function.Function
type Function[T any, R any] func(T) R

// Predicate
// See: java.util.function.Predicate
type Predicate[T any] func(T) bool

// Supplier
// See: java.util.function.Supplier
type Supplier[T any] func() T

// UnaryOperator
// See: java.util.function.UnaryOperator
type UnaryOperator[T any] Function[T, T]
