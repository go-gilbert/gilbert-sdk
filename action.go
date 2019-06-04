package sdk

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"sync"
	"time"
)

// Actions is actions map
//
// Key is an action name and value is action constructor
type Actions map[string]HandlerFactory

// ActionParams is action params container
type ActionParams map[string]interface{}

// Unmarshal extracts action params into provided structure
func (p ActionParams) Unmarshal(dest interface{}) error {
	if err := mapstructure.Decode(p, dest); err != nil {
		return fmt.Errorf("failed to unmarshal action params, %s", err)
	}

	return nil
}

// HandlerFactory is action handler constructor
type HandlerFactory func(ScopeAccessor, ActionParams) (ActionHandler, error)

// ActionHandler represents Gilbert's action handler
type ActionHandler interface {
	// Call calls an action handler
	Call(JobContextAccessor, JobRunner) error

	// Cancel aborts action handler execution
	Cancel(JobContextAccessor) error
}

// JobRunner is the the interface that represents a current job caller.
type JobRunner interface {
	// ActionByName returns action handler constructor
	ActionByName(actionName string) (p HandlerFactory, err error)

	// RunJob starts job in separate goroutine.
	//
	// Use ctx.Errors() to track job result and ctx.Cancel() to cancel job execution.
	RunJob(j Job, ctx JobContextAccessor)

	// RunTask runs task by name
	//
	// Use ctx.Errors() to track job result and ctx.Cancel() to cancel job execution.
	RunTask(taskName string, ctx JobContextAccessor, scope ScopeAccessor) error
}

// JobContextAccessor provides access to job run context used store job state and communicate between task runner and job
type JobContextAccessor interface {
	// Log provides logger for current job context
	Log() Logger

	// IsAlive checks if context was not finished
	IsAlive() bool

	// IsChild checks if context is child context
	IsChild() bool

	// Context returns Go context instance assigned to the current job context
	Context() context.Context

	// Errors returns channel with job execution result
	Errors() chan error

	// SetErrorChannel sets custom error report channel
	SetErrorChannel(chan error)

	// Vars returns a set of variables attached to this context
	Vars() Vars

	// SetVars sets context variables
	SetVars(Vars)

	// SetWaitGroup sets wait group instance for current job
	//
	// This value will be used later to call wg.Done() when job was finished.
	SetWaitGroup(wg *sync.WaitGroup)

	// ForkContext creates a job context copy, but creates a separate sub-logger
	ForkContext() JobContextAccessor

	// ChildContext creates a new child job context with separate Error channel and context
	ChildContext() JobContextAccessor

	// Timeout adds timeout to the context
	Timeout(timeout time.Duration)

	// Success reports successful result.
	//
	// Alias to Result(nil)
	Success()

	// Result reports job result and finished the context
	Result(err error)

	// Cancel cancels the context and stops all jobs related to the context
	Cancel()
}
