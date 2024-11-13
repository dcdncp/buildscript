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
	v, stt := CheckArgsCount(env, len(args), 1, false)
	if stt.IsNotOkay() {
		return v, stt
	}
	v, stt = CheckType(env, args[0], StringType)
	if stt.IsNotOkay() {
		return v, stt
	}
	p := v.(*String).Value
	source := env.Global.Source
	file := env.Global.SourceFile
	if strings.HasPrefix(p, "std/") {
		p :=  strings.TrimPrefix(p, "std/")
		e, exists := Modules[p]
		if exists {
			env.Global.Source = source
			env.Global.SourceFile = file
			return NewModule(e), state.Ok
		}
		cp, err := os.Getwd()
		if err != nil {
			panic(err)
		} 
		p = path.Join(cp, "std", p) + ".bs"
		if _, err := os.Stat(p); err == nil {
			m, exists := LoadedModules[p]
			if exists {
				env.Global.Source = source
				env.Global.SourceFile = file
				return m, state.Ok
			}
			v, stt := env.Global.EvalFile(p)
			if stt.IsOkay() {
				LoadedModules[p] = v
			}
			env.Global.Source = source
			env.Global.SourceFile = file
			return v, stt
		}
	}
	if !path.IsAbs(p) {
		p = path.Join(path.Dir(env.Global.SourceFile), p)
	}
	m, exists := LoadedModules[p]
	if exists {
		env.Global.Source = source
		env.Global.SourceFile = file
		return m, state.Ok
	}
	v, stt = env.Global.EvalFile(p)
	if stt.IsOkay() {
		LoadedModules[p] = v
	}
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
