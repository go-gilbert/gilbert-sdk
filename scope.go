package sdk

// Scope contains a set of globals and variables related to specific job
type Scope interface {

	// AppendVariables appends local variables to the context
	AppendVariables(vars Vars) Scope

	// AppendGlobals appends global variables to the context
	AppendGlobals(vars Vars) Scope

	// Global returns a global variable value by it's name
	Global(varName string) (out string, ok bool)

	// Var returns a variable value by it's name
	Var(varName string) (isLocal bool, out string, ok bool)

	// ExpandVariables expands an expression stored inside a passed string
	ExpandVariables(expression string) (string, error)

	// Scan does the same as ExpandVariables but with multiple variables and updates the value in pointer with expanded value
	//
	// Useful for bulk mapping of struct fields
	Scan(values ...*string) (err error)

	// Environ gets list of OS environment variables with globals
	Environ() (env []string)
}
