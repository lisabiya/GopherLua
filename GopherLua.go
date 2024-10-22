package GopherLua

import (
	lua "github.com/yuin/gopher-lua"
	"math"
)

type Lua struct {
	State  *lua.LState
	module []*Module
}

//拓展模块需要实现的接口
type Module interface {
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

func (instance *Lua) Register(module ...Module) {
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

//获取栈中指定返回值，并pop出去，节约内存
func (instance *Lua) GetAndPop(idx int) lua.LValue {
	var value = instance.State.Get(idx)
	instance.State.Pop(int(math.Abs(float64(idx))))
	return value
}
