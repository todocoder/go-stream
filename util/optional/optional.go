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

package optional

import "github.com/zerune/go-core/util/function"

type Optional[T any] struct {
	value *T
}

func newOptional[T any](value *T) Optional[T] {
	return Optional[T]{value: value}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

func Of[T any](value T) Optional[T] {
	return newOptional(&value)
}

func OfNullable[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return newOptional(value)
}

func (o Optional[T]) Get() *T {
	if o.value == nil {
		return new(T)
	}
	return o.value
}

func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o Optional[T]) IsEmpty() bool {
	return o.value == nil
}

func (o Optional[T]) IfPresent(action function.Consumer[*T]) {
	if o.value != nil {
		action(o.value)
	}
}

func (o Optional[T]) IfPresentOrElse(action function.Consumer[*T], emptyAction function.Runnable) {
	if o.value != nil {
		action(o.value)
	} else {
		emptyAction()
	}
}

func (o Optional[T]) Filter(predicate function.Predicate[*T]) Optional[T] {
	if !o.IsPresent() {
		return o
	} else {
		if predicate(o.value) {
			return o
		}
		return Empty[T]()
	}
}

func (o Optional[T]) Map(mapper function.Function[*T, any]) Optional[any] {
	if !o.IsPresent() {
		return Empty[any]()
	} else {
		r := mapper(o.value)
		return OfNullable(&r)
	}
}

func (o Optional[T]) FlatMap(mapper function.Function[*T, Optional[any]]) Optional[any] {
	if !o.IsPresent() {
		return Empty[any]()
	} else {
		return mapper(o.value)
	}
}

func (o Optional[T]) Or(supplier function.Supplier[Optional[T]]) Optional[T] {
	if o.IsPresent() {
		return o
	} else {
		return supplier()
	}
}

func (o Optional[T]) OrElse(other *T) *T {
	if o.value != nil {
		return o.value
	} else {
		return other
	}
}

func (o Optional[T]) OrElseGet(supplier function.Supplier[*T]) *T {
	if o.value != nil {
		return o.value
	} else {
		return supplier()
	}
}
