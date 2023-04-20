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

package errors

import err "errors"

func New(message string) error {
	return err.New(message)
}

type RuntimeError struct {
	message string
}

func NewRuntimeError(message string) *RuntimeError {
	return &RuntimeError{message: message}
}

func (e *RuntimeError) Error() string {
	return e.message
}

// IllegalArgumentError Thrown to indicate that a method has been passed an illegal or inappropriate argument.
type IllegalArgumentError struct {
	RuntimeError
}

func NewIllegalArgumentError(message string) *IllegalArgumentError {
	return &IllegalArgumentError{RuntimeError: *NewRuntimeError(message)}
}

// UnsupportedOperationError Thrown to indicate that the requested operation is not supported.
type UnsupportedOperationError struct {
	RuntimeError
}

func NewUnsupportedOperationError(message string) *UnsupportedOperationError {
	return &UnsupportedOperationError{RuntimeError: *NewRuntimeError(message)}
}
