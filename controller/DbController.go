package controller

import (
	"GopherLua/db_module"
	"GopherLua/goTool"
	"github.com/gin-gonic/gin"
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

func LoadLuaModule(c *gin.Context) {
	luaContext := getDefaultGinStatus(c)
	db_module.RegisterOrmDbType(luaContext)
	defer luaContext.Close()
	err := luaContext.DoFile("lua/run.lua")
	if err != nil {
		c.JSON(http.StatusOK, formatError(err))
		return
	}
	if err := luaContext.CallByParam(lua.P{
		Fn:      luaContext.GetGlobal("initParams"),
		NRet:    1,
		Protect: true,
	}); err != nil {
		c.JSON(http.StatusOK, formatError(err))
		return
	}
	ret := luaContext.Get(1) // returned value

	c.JSON(http.StatusOK, formatSuccess(goTool.TransLuaValue2Map(ret)))
}

func getDefaultGinStatus(c *gin.Context) *lua.LState {
	L := lua.NewState()
	var getParams = L.NewFunction(func(state *lua.LState) int {
		var key = state.ToString(-1)
		var value = c.Query(key)
		L.Push(lua.LString(value))
		return 1
	})
	L.SetGlobal("getParams", getParams)
	var postParams = L.NewFunction(func(state *lua.LState) int {
		var key = state.ToString(-1)
		var value = c.PostForm(key)
		L.Push(lua.LString(value))
		return 1
	})
	L.SetGlobal("postParams", postParams)
	return L
}
