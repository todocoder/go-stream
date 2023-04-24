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
	"errors"
	"fmt"
	"reflect"
)

// Try
// see: https://www.cnblogs.com/beiluowuzheng/p/10263724.html
func Try(f func() error) Catch {
	c := &catch{}
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					c.err = err
				} else {
					c.err = errors.New(fmt.Sprint(err))
				}
			}
		}()
		if err := f(); err != nil {
			c.err = err
		}
	}()
	return c
}

type Catch interface {
	Catch(error, func(error)) Catch
	CatchAll(func(error)) Finally
	Finally(...func())
}

type Finally interface {
	Finally(...func())
}

type catch struct {
	err    error
	caught bool
}

func (c *catch) required() bool {
	if c.caught || c.err == nil {
		return false
	}
	return true
}

func (c *catch) Catch(err error, handler func(error)) Catch {
	if !c.required() {
		return c
	}
	// 如果传入的error类型和发生异常的类型一致，则执行异常处理器，并将caught修改为true代表已捕捉异常
	if reflect.TypeOf(err) == reflect.TypeOf(c.err) {
		handler(c.err)
		c.caught = true
	}
	return c
}

func (c *catch) CatchAll(handler func(error)) Finally {
	if !c.required() {
		return c
	}
	handler(c.err)
	c.caught = true
	return c
}

func (c *catch) Finally(handlers ...func()) {
	defer func() {
		// 遍历处理器，并在Finally函数执行完毕之后执行
		for _, handler := range handlers {
			handler()
		}
	}()
	// 如果异常不为空，且未捕捉异常，则抛出异常
	if c.err != nil && !c.caught {
		panic(c.err)
	}
}
