package starlarkgame

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func NewInterpreter(predeclared starlark.StringDict) *Interpreter {
	predeclared["struct"] = starlark.NewBuiltin("struct", starlarkstruct.Make)
	return &Interpreter{predeclared: predeclared}
}

func (i *Interpreter) ExecFile(filename string) error {
	thread := &starlark.Thread{
		Name: "ExecFile",
		Load: i.load,
	}
	globals, err := starlark.ExecFile(thread, filename, nil, i.predeclared)
	if err != nil {
		return err
	}
	i.globals = globals
	return err
}

type Interpreter struct {
	predeclared starlark.StringDict
	globals     starlark.StringDict
}

// Call executes fn with the given positional args and keyword kwargs.
//
// Returns the result of the call, which may be starlark.None if fnName returns no value.
func (i *Interpreter) Call(fnName string, args, kwargs starlark.Tuple) (starlark.Value, error) {
	thread := &starlark.Thread{
		Name: fnName,
		Load: i.load,
	}
	fn, ok := i.globals[fnName]
	if !ok {
		return nil, fmt.Errorf("no such method: %q", fnName)
	}
	return starlark.Call(thread, fn, args, nil)
}

// load implements the 'load' operation.
//
// If the module name begins with '//' then the module path is treated as being relative
// to the directory of the original filename passed to ExecFile. Otherwise it is
// treated as being relative to the module that is currently executing.
func (i *Interpreter) load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	var dir string
	if strings.HasPrefix(module, "//") {
		module = module[len("//"):]
		dir = filepath.Dir(thread.CallFrame(thread.CallStackDepth() - 1).Pos.Filename())
	} else {
		dir = filepath.Dir(thread.CallFrame(0).Pos.Filename())
	}
	filename := filepath.Join(dir, module)
	return starlark.ExecFile(thread, filename, nil, i.predeclared)
}
