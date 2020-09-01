package httpRequest

import (
	"GopherLua/goTool"
	"fmt"
	"github.com/parnurzeal/gorequest"
	lua "github.com/yuin/gopher-lua"
)

const metatableName = "request_metatable"

func RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(metatableName)
	L.SetGlobal("httpRequest", mt)
	//初始化实例
	L.SetField(mt, "new", L.NewFunction(newObject))
	//方法
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), gormMethods))
}

func newObject(L *lua.LState) int {
	ud := L.NewUserData()
	L.SetMetatable(ud, L.GetTypeMetatable(metatableName))
	L.Push(ud)
	return 1
}

//链式调用需要来回传输实例，暂时以参数形式一次调用
var gormMethods = map[string]lua.LGFunction{
	"postForm": postForm,
	"getQuery": getQuery,
}

func postForm(L *lua.LState) int {
	var request = L.CheckTable(2)
	var requestMap, ok = goTool.TransLuaValue2Map(request).(map[string]interface{})
	if ok {
		_, body, errs := gorequest.New().
			Post(requestMap["url"].(string)).
			Type(requestMap["type"].(string)).
			SendMap(requestMap["params"]).End()
		//
		if len(errs) > 0 {
			var errStr = ""
			for _, err := range errs {
				errStr = errStr + err.Error() + "\n"
			}
			L.Push(lua.LNumber(1))
			L.Push(lua.LString(errStr))
			return 2
		}
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(body))
		return 2
	} else {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString("参数转map失败"))
		return 2

	}
}

func getQuery(L *lua.LState) int {
	var request = L.CheckTable(2)
	var requestMap, ok = goTool.TransLuaValue2Map(request).(map[string]interface{})
	if ok {
		req, body, errs := gorequest.New().
			Get(requestMap["url"].(string)).
			Query(requestMap["query"]).End()
		if len(errs) > 0 {
			var errStr = ""
			for _, err := range errs {
				errStr = errStr + err.Error() + "\n"
			}
			L.Push(lua.LNumber(1))
			println(errStr)
			L.Push(lua.LString(errStr))
			return 2
		}
		fmt.Println(req.Request.RequestURI)
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(body))
		return 2
	} else {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString("参数转map失败"))
		return 2

	}
}
