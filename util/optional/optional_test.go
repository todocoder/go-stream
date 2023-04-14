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

import (
	"reflect"
	"testing"
)

type testItem struct {
	id    int
	value string
}

func TestEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		want *Optional[T]
	}
	tests := []testCase[string]{
		{name: "empty", want: &Optional[string]{}},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Empty[string](); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Empty() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestOptional_Get(t *testing.T) {
	a := Empty[string]()
	if a.IsEmpty() {
		t.Log("empty")
	}
	a = Of[string]("1")
	got := a.Get()
	if !reflect.DeepEqual(got, "1") {
		t.Errorf("Get() got = %v, want %v", got, "1")
	}
}

func TestOptional_Map(t *testing.T) {
	item := testItem{id: 1, value: "item1"}
	o := Of[testItem](item).Map(func(item testItem) any {
		return item.value
	})
	value := o.Get()
	if v, ok := value.(string); ok {
		t.Logf("got = %v", *o.value)
	} else {
		t.Errorf("got = %v, want %v", v, "item1")
	}
}

func TestOptional_FlatMap(t *testing.T) {
	item := testItem{id: 1, value: "item1"}
	o := Of[testItem](item).FlatMap(func(item testItem) *Optional[any] {
		return Of[any](item.value)
	})
	value := o.Get()
	if v, ok := value.(string); ok {
		t.Logf("got = %v", *o.value)
	} else {
		t.Errorf("got = %v, want %v", v, "item1")
	}
}
