package module_http

import (
	"github.com/lisabiya/GopherLua/goTool"
	"github.com/parnurzeal/gorequest"
	lua "github.com/yuin/gopher-lua"
)

//网络请求模块
type ModuleHttp struct {
}

const metatableName = "request_metatable"

func (http ModuleHttp) RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(metatableName)
	L.SetGlobal("httpRequest", mt)
	//初始化实例
	L.SetField(mt, "new", L.NewFunction(newObject(&http)))
	//方法
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))
}

func (http ModuleHttp) Close() {

}

//**********元表拓展给对象的方法***********

func newObject(http *ModuleHttp) lua.LGFunction {
	return func(state *lua.LState) int {
		ud := state.NewUserData()
		ud.Value = http
		state.SetMetatable(ud, state.GetTypeMetatable(metatableName))
		state.Push(ud)
		return 1
	}
}

//**********元表拓展给对象的方法***********

var methods = map[string]lua.LGFunction{
	"End": end,
}

func end(L *lua.LState) int {
	var request = L.CheckTable(2)
	var requestMap, ok = goTool.TransLuaValue2Map(request).(map[string]interface{})
	if ok {
		var request = gorequest.New()
		if requestMap["get"] != nil {
			request.Get(requestMap["get"].(string))
		}
		if requestMap["post"] != nil {
			request.Post(requestMap["post"].(string))
		}

		for option, value := range requestMap {
			println(option)
			switch option {
			case "query":
				request.Query(value)
				break
			case "type":
				request.Type(value.(string))
				break
			case "send":
				request.Send(value)
				break
			case "set":
				var headers = value.(map[string]interface{})
				for k, v := range headers {
					request.Set(k, goTool.FormatString(v))
				}
				break
			}
		}
		res, body, errs := request.End()
		println(res.Request.Header)
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

func checkObject(L *lua.LState) *ModuleHttp {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*ModuleHttp); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}
