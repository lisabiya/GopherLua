package GopherLua

import lua "github.com/yuin/gopher-lua"

type Lua struct {
	State  *lua.LState
	module []*InterfaceModule
}

type InterfaceModule interface {
	RegisterType(L *lua.LState)
	Close()
}

func NewState() *Lua {
	return &Lua{
		State: lua.NewState(),
	}
}

func (instance *Lua) DoString(source string) error {
	return instance.State.DoString(source)
}

func (instance *Lua) DoFile(path string) error {
	return instance.State.DoFile(path)
}

func (instance *Lua) ExecuteFunc(funcName string, returnParamsCount int, args ...lua.LValue) error {
	return instance.State.CallByParam(lua.P{
		Fn:      instance.State.GetGlobal(funcName),
		NRet:    returnParamsCount,
		Protect: false,
	}, args...)
}

func (instance *Lua) Register(module ...InterfaceModule) {
	for _, interfaceModule := range module {
		instance.module = append(instance.module, &interfaceModule)
		interfaceModule.RegisterType(instance.State)
	}
}

func (instance *Lua) Close() {
	for i := range instance.module {
		var module = instance.module[i]
		(*module).Close()
	}
}
