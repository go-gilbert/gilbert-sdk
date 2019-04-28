package sdk

import (
	"context"
	"sync"
	"time"
)

// PluginParams is raw plugin params
type PluginParams map[string]interface{}

// PluginFactory is plugin constructor
type PluginFactory func(ScopeAccessor, PluginParams, Logger) (Plugin, error)

// Plugin represents Gilbert's plugin
type Plugin interface {
	// Call calls a plugin
	Call(JobContextAccessor, JobRunner) error

	// Cancel stops plugin execution
	Cancel(JobContextAccessor) error
}

// JobRunner is the the interface that represents a current job caller.
type JobRunner interface {
	// PluginByName returns plugin constructor
	PluginByName(pluginName string) (p PluginFactory, err error)

	// RunJob starts job in separate goroutine.
	//
	// Use ctx.Error channel to track job result and ctx.Cancel() to cancel it.
	RunJob(j Job, ctx JobContextAccessor)
}

// JobContextAccessor provides access to job run context used store job state and communicate between task runner and job
type JobContextAccessor interface {

	// IsAlive checks if context was not finished
	IsAlive() bool

	// IsChild checks if context is child context
	IsChild() bool

	// Context returns Go context instance assigned to the current job context
	Context() context.Context

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
