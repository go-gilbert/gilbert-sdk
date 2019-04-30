package sdk

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"sync"
	"time"
)

// ActionParams is plugin params container
type ActionParams map[string]interface{}

// Unmarshal extracts plugin params into provided structure
func (p ActionParams) Unmarshal(dest interface{}) error {
	if err := mapstructure.Decode(p, dest); err != nil {
		return fmt.Errorf("failed to unmarshal plugin params, %s", err)
	}

	return nil
}

// ActionFactory is action handler constructor
type ActionFactory func(ScopeAccessor, ActionParams, Logger) (ActionHandler, error)

// ActionHandler represents Gilbert's action handler
type ActionHandler interface {
	// Call calls a plugin
	Call(JobContextAccessor, JobRunner) error

	// Cancel stops plugin execution
	Cancel(JobContextAccessor) error
}

// JobRunner is the the interface that represents a current job caller.
type JobRunner interface {
	// ActionByName returns action handler constructor
	ActionByName(actionName string) (p ActionFactory, err error)

	// RunJob starts job in separate goroutine.
	//
	// Use ctx.Error channel to track job result and ctx.Cancel() to cancel it.
	RunJob(j Job, ctx JobContextAccessor)
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

	// Errors returns job errors channel
	Errors() chan error

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

	// Cancel cancels the context and stops all jobs used by this context
	Cancel()
}
