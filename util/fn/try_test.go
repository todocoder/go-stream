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

package fn

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/zerune/go-core/lang/errors"
)

func TestTry(t *testing.T) {
	Try(func() error {
		n := rand.Intn(9)
		fmt.Println(n)
		if n > 5 {
			panic(errors.New("Greater than 5"))
		}
		r := n%2 == 0
		if r {
			return errors.NewUnsupportedOperationError("UnsupportedOperationError")
		}
		return errors.NewIllegalArgumentError("IllegalArgumentError")
	}).Catch(new(errors.IllegalArgumentError), func(err error) {
		fmt.Println("Catch IllegalArgumentError")
	}).Catch(new(errors.UnsupportedOperationError), func(err error) {
		fmt.Println("Catch UnsupportedOperationError")
	}).CatchAll(func(err error) {
		fmt.Println(fmt.Sprint(err))
	}).Finally(func() {
		fmt.Println("Finally")
	})
}
