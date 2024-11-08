package fs

import (
	"bscript/runtime/state"
	"bscript/runtime/std"
	"bscript/runtime/value"
	"os"
	"path"
)

func readFile(env *value.Env, args []value.Value) (value.Value, state.State) {
	p := args[0].(*std.String).Value
	if !path.IsAbs(p) {
		p = path.Join(path.Join(path.Dir(env.Global.SourceFile), p))
	}
	content, err := os.ReadFile(p)
	if err != nil {
		return std.ThrowException(env, err.Error())
	}
	return std.NewString(string(content)), state.Ok
}

func writeFile(env *value.Env, args []value.Value) (value.Value, state.State) {
	p := args[0].(*std.String).Value
	if !path.IsAbs(p) {
		p = path.Join(path.Join(path.Dir(env.Global.SourceFile), p))
	}
	content := args[1].(*std.String).Value
	err := os.WriteFile(p, []byte(content), os.ModePerm)
	if err != nil {
		return std.ThrowException(env, err.Error())
	}
	return nil, state.Ok
}
