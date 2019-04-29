package sdk

// ProjectEnvironment contains information about project environment
type ProjectEnvironment struct {
	ProjectDirectory string
}

// ScopeAccessor contains a set of globals and variables related to specific job
type ScopeAccessor interface {

	// AppendVariables appends local variables to the context
	AppendVariables(vars Vars) ScopeAccessor

	// AppendGlobals appends global variables to the context
	AppendGlobals(vars Vars) ScopeAccessor

	// Global returns a global variable value by it's name
	Global(varName string) (out string, ok bool)

	// Var returns a variable value by it's name
	Var(varName string) (isLocal bool, out string, ok bool)

	// Variables returns all declared local variables
	Vars() Vars

	// ExpandVariables expands an expression stored inside a passed string
	ExpandVariables(expression string) (string, error)

	// Scan does the same as ExpandVariables but with multiple variables and updates the value in pointer with expanded value
	//
	// Useful for bulk mapping of struct fields
	Scan(values ...*string) (err error)

	// Environment returns information about project environment
	Environment() ProjectEnvironment

	// Environ gets list of OS environment variables with globals
	Environ() (env []string)
}
