package std

import (
	"bscript/runtime/state"
	"bscript/runtime/value"
	"fmt"
	"os"
	"path"
	"strings"
)

func importModule(env *value.Env, args []value.Value) (value.Value, state.State) {
	if len(args) != 1 {
		return NewString("incorrect count of arguments"), state.Error
	}
	source := env.Global.Source
	file := env.Global.SourceFile
	v := args[0]
	if v.Type() != StringType {
		return NewString("function 'import' requires string argument"), state.Error
	}
	p := v.(*String).Value
	if strings.HasPrefix(p, "std/") {
		p :=  strings.TrimPrefix(p, "std/")
		e, exists := Modules[p]
		if exists {
			return NewModule(e), state.Ok
		}
		cp, err := os.Getwd()
		if err != nil {
			panic(err)
		} 
		p = path.Join(cp, "std", p) + ".bs"
		if _, err := os.Stat(p); err == nil {
			v, stt := env.Global.EvalFile(p)
			env.Global.Source = source
			env.Global.SourceFile = file
			return v, stt
		}
	}
	if !path.IsAbs(p) {
		p = path.Join(path.Dir(env.Global.SourceFile), p)
	}
	v, stt := env.Global.EvalFile(p)
	env.Global.Source = source
	env.Global.SourceFile = file
	return v, stt
}

func message(env *value.Env, args []value.Value) (value.Value, state.State) {
	strs := make([]string, 0)
	for _, v := range args {
		s, stt := ToString(env, v)
		if stt.IsNotOkay() {
			return s, stt
		}
		strs = append(strs, s.(*String).Value)
	}
	fmt.Println(strings.Join(strs, " "))
	return nil, state.Ok
}

func typeOf(env *value.Env, args []value.Value) (value.Value, state.State) {
	v := args[0]
	return v.Type(), state.Ok
}
