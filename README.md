# GopherLua 🚜
lua为go增加动态化能力，go为lua提供功能拓展

>项目依托[gopher-lua](https://github.com/yuin/gopher-lua)`go平台的lua解释器` 
><br>进行拓展和封装，主要是提供一个方向/思路，并尝试实现 

## go拓展库
- [x] 数据库连接查询库`module_db`
- [x] 网络请求库`module_http`

### 安装(Installation)
```go
go get github.com/lisabiya/GopherLua
```

### 简单示例--http请求(simple Example httpRequest)
```go
import (
	"fmt"
	"github.com/lisabiya/GopherLua"
	"github.com/lisabiya/GopherLua/module_http"
)

func main() {
	gopherLua := GopherLua.NewState()
	gopherLua.Register(module_http.ModuleHttp{})
	err := gopherLua.DoString(
		`
    --引用声明模块
	local http = httpRequest.new()
    local code, response =  http:End({
        get = "https://www.wanandroid.com/hotkey/json",
        query = "nihao",
    })
    --调用函数  
    print(code,response)
`)
	if err != nil {
		fmt.Println(err.Error())
	}
}
```
### 性能(performance)
- 在postman测试脚本中，mysql请求数据列表20ms,1000请求下，内存基本稳定10M以内无变化
- 当然这只是浅略的测试，目前正打算逐步加入正式项目中实践。（队友听了想打人😀）


<br><br>

---
 
### 添加自定义拓展(Add custom extension)
> 主要是提供一种思路，需要优化改进的地方还有很多 
 
- 参考 `module_http,module_db` 

#### 示例一 简单声明拓展函数

- go声明函数

```go

const metatableName = "request_metatable"

func RegisterFuction(L *lua.LState) {
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

func main() {
	gopherLua := GopherLua.NewState()
	RegisterFuction(gopherLua.State)
	err := gopherLua.DoFile("test/test.lua")//test.lua 为项目所在文件位置
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

- lua中调用 `test.lua`
```lua
--直接引用声明模块
local code, response = httprequest.get(
        { url = "https://www.wanandroid.com/hotkey/json" })
--调用函数
print(code, response)
```

#### 示例二 声明模块

- go声明函数

```go
const metatableName = "request_metatable"

type CustomModule struct {
}

func (CustomModule) RegisterType(L *lua.LState) {
	mt := L.NewTypeMetatable(metatableName)
	//声明全局对象
	L.SetGlobal("httprequest", mt)
	//添加拓展函数
	L.SetField(mt, "get", L.NewFunction(getSimple))
}

func (CustomModule) Close() {

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

func main() {
	gopherLua := GopherLua.NewState()
	gopherLua.Register(CustomModule{})
	err := gopherLua.DoFile("test/test.lua")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = gopherLua.ExecuteFunc("TestHttp", 1, lua.LString("测试参数"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(goTool.TransLuaValue2Map(gopherLua.State.Get(-1)))
}
```
- lua中调用 `test.lua`
```lua
--直接引用声明模块
function TestHttp(params)
    local code, response = httprequest.get(
            { url = "https://www.wanandroid.com/hotkey/json" })
    --调用函数
    return { code = code, response = response, params = params }
end
```




