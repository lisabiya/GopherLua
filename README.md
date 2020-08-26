# GopherLua 🚜
lua为go增加动态化能力，go为lua提供功能拓展

> 项目依托[gopher-lua](https://github.com/yuin/gopher-lua)`go平台的lua解释器` 对lua进行拓展

## go拓展
- [x] 数据库连接查询库`db_module`
- [ ] 网络请求库

> 简单示例-数据库连接查询库
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


>### 添加拓展-`具体实例参考 db_module/example.lua`
- go声明元表函数

```go
    mt := L.NewTypeMetatable("db_metatable")
    //声明全局模块名
    L.SetGlobal("db_module", mt)
    //创建实例，添加元表
    L.SetField(mt, "new", L.NewFunction(
        func(state *lua.LState) int {
		var DbName = L.CheckString(1)
		ormDb := &OrmDB{DbName: DbName}
		ud := L.NewUserData()
		ud.Value = ormDb
		L.SetMetatable(ud, L.GetTypeMetatable("db_metatable"))
		L.Push(ud)
		return 1
	}))
    //添加元表函数
    L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), 
        map[string]lua.LGFunction{
			"getName": func(state *lua.LState) int {
                   		p := checkDb(state)
                   		state.Push(lua.LString(p.DbName))
                   		return 1
                   	}}))
```

- lua中调用
```lua
  --直接引用声明模块
  local obj = db_module.new("DbName")
  --调用函数  
  print(obj:getName())
```






