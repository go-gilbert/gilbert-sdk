package sdk

import (
	"time"
)

// JobExecType represents job type
type JobExecType uint8

const (
	// ExecEmpty means job has no execution type
	ExecEmpty JobExecType = iota

	// ExecAction means that job execute action
	ExecAction

	// ExecMixin means that job based on mixin
	ExecMixin

	// ExecTask means that job should execute another task
	ExecTask
)

// Period is job period in milliseconds
type Period uint

// ToDuration returns value in milliseconds for time.Duration
func (d Period) ToDuration() time.Duration {
	return time.Duration(d) * time.Millisecond
}

// Job represents a single step in task
type Job struct {
	// Condition is shell command that should be successful to run specified job
	Condition string `yaml:"if,omitempty" mapstructure:"if"`

	// Description is job description
	Description string `yaml:"description,omitempty" mapstructure:"description"`

	// TaskName refers to task that should be run.
	TaskName string `yaml:"task,omitempty" mapstructure:"run"`

	// ActionName is action to run.
	ActionName string `yaml:"action,omitempty" mapstructure:"action"`

	// MixinName is mixin to run.
	//
	// Cannot present with ActionName at the same time
	MixinName string `yaml:"mixin,omitempty" mapstructure:"mixin"`

	// Async means that job should be run asynchronously
	Async bool `yaml:"async,omitempty" mapstructure:"async"`

	// Delay before task start in milliseconds
	Delay Period `yaml:"delay,omitempty" mapstructure:"delay"`

	// Period is a time quota for job
	Deadline Period `yaml:"deadline,omitempty" mapstructure:"deadline"`

	// Vars is a set of variables defined for this job.
	Vars Vars `yaml:"vars,omitempty" mapstructure:"vars"`

	// Params is a set of arguments for the job.
	Params ActionParams `yaml:"params,omitempty" mapstructure:"params"`
}

// HasDescription checks if description is available
func (j *Job) HasDescription() bool {
	return j.Description != ""
}

// FormatDescription returns formatted description string
func (j *Job) FormatDescription() string {
	if j.Description != "" {
		return j.Description
	}

	// If description is empty, return used mixin or action name if available
	for _, v := range []string{j.ActionName, j.TaskName, j.MixinName} {
		if v != "" {
			return v
		}
	}

	return ""
}

// Type returns job execution type
//
// If job has no 'action', 'task' or 'mixin' declaration, ExecEmpty will be returned
func (j *Job) Type() JobExecType {
	if j.ActionName != "" {
		return ExecAction
	}

	if j.MixinName != "" {
		return ExecMixin
	}

	if j.TaskName != "" {
		return ExecTask
	}

	return ExecEmpty
}
