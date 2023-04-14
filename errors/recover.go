/*
 * See: zeromicro/go-zero/core/rescue/recover.go
 */

package errors

import "log"

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		log.Println(p)
	}
}
