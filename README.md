# GopherLua 🚜
lua为go增加动态化能力，go为lua提供功能拓展

> 项目依托[gopher-lua](https://github.com/yuin/gopher-lua)`go平台的lua解释器` 对lua进行拓展

## go拓展
- [x] 数据库连接查询库`db_module`
- [x] 网络请求库`httpRequest`


> 简单示例-数据库请求
```lua
local luaDbSqLite = require('db_module.db_module')
local ormDb = luaDbSqLite.new("t_salary")
local Builder = require "db_module.LuaQuB"

function example.getList()
    local object = Builder.new():select("*"):from("t_salary") :limit(10, 0)
    local code, tables = ormDb:Raw(tostring(object))
    print(ormDb:Tag(), #tables)
    return code, { count = #tables, list = tables }
end
```


>### 添加拓展-`具体实例参考 /httprequest/httpExample.lua`
- go声明元表函数

```go
const metatableName = "request_metatable"

func RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(metatableName)
    //声明全局对象
	L.SetGlobal("httprequest", mt)
    //添加拓展函数
	L.SetField(mt, "get", L.NewFunction(getSimple))
}

func getSimple(L *lua.LState) int {
	var request = L.CheckTable(1)
	var requestMap, ok = goTool.TransLuaValue2Map(request).(map[string]interface{})
	if ok {
		_, body, _ := gorequest.New().
			Get(requestMap["url"].(string)).
			Query(requestMap["query"]).End()
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(body))
		return 2
	} else {
		L.Push(lua.LNumber(1))
		L.Push(lua.LString("参数转map失败"))
		return 2

	}
}
```

- lua中调用
```lua
  --直接引用声明模块
    local code, response = httprequest:get(
            { url = "https://www.wanandroid.com/hotkey/json" })
  --调用函数  
    print(code,response)
```






