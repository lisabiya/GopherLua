# GopherLua 🚜
lua为go增加动态化能力，go为lua提供功能拓展

>项目依托[gopher-lua](https://github.com/yuin/gopher-lua)`go平台的lua解释器` 对lua进行拓展

## go拓展库
- [x] 数据库连接查询库`module_db`
- [x] 网络请求库`module_http`

### 安装--Installation
```go
go get github.com/lisabiya/GopherLua
```

### 简单示例-http请求 (simple Example httpRequest)
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

<br><br>

### 添加自定义拓展
> 主要是提供一种思路，需要优化改进的地方还有很多 
 
- 参考 `module_http,module_db` 

- #### 示例

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







