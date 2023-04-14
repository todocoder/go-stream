/*
 * See: zeromicro/go-zero/core/threading/routinegroup.go
 */

package routine

import (
	"sync"

	"github.com/zerune/go-core/lang"
)

type WaitGroup interface {
	Run(fn func())
	RunSafe(fn func())
	Wait()
}

// A RoutineGroup is used to group goroutines together and all wait all goroutines to be done.
type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

// NewRoutineGroup returns a RoutineGroup.
func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}

// Run runs the given fn in RoutineGroup.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// RunSafe runs the given fn in RoutineGroup, and avoid panics.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	})
}

// Wait waits all running functions to be done.
func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}

// A LimitedGroup is used to group goroutines together and all wait all goroutines to be done.
type LimitedGroup struct {
	waitGroup sync.WaitGroup
	pool      chan lang.PlaceholderType
}

// NewLimitedGroup returns a RoutineGroup.
func NewLimitedGroup(workers ...int) *LimitedGroup {
	size := 1
	for _, worker := range workers {
		if worker > size {
			size = worker
		}
	}
	return &LimitedGroup{
		pool: make(chan lang.PlaceholderType, size),
	}
}

// Run runs the given fn in RoutineGroup.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *LimitedGroup) Run(fn func()) {
	g.pool <- lang.Placeholder
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		<-g.pool
		fn()
	}()
}

// RunSafe runs the given fn in RoutineGroup, and avoid panics.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *LimitedGroup) RunSafe(fn func()) {
	g.pool <- lang.Placeholder
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		<-g.pool
		fn()
	})
}

// Wait waits all running functions to be done.
func (g *LimitedGroup) Wait() {
	g.waitGroup.Wait()
	close(g.pool)
}
