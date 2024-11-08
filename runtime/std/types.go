package std

import (
	"bscript/runtime/state"
	"bscript/runtime/symbol"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"bscript/tools"
	"fmt"
)

func typeString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*types.Type)
	return NewString(fmt.Sprintf("<type %s>", self.Name)), state.Ok
}

func typeCall(env *value.Env, args []value.Value) (value.Value, state.State) {
	t := args[0].(*types.Type)
	self := t.New()
	f, exists := t.GetField(symbol.Init)
	if exists {
		self.Unlock()
		v, state := Call(env, f, tools.AppendFront(args[1:], self)...)
		if state.IsNotOkay() {
			self.Lock()
			return v, state
		}
		self.Unlock()
	}
	return self, state.Ok
}
