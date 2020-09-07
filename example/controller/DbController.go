package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lisabiya/GopherLua"
	"github.com/lisabiya/GopherLua/goTool"
	"github.com/lisabiya/GopherLua/module_db"
	"github.com/lisabiya/GopherLua/module_http"
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

var gopherLua *GopherLua.Lua

func init() {
	gopherLua = GopherLua.NewState()
	//引入模块
	gopherLua.Register(module_db.ModuleDb{}, module_http.ModuleHttp{})
}

func LoadLuaModule(c *gin.Context) {
	//自定义功能
	setDefaultGinStatus(c, gopherLua.State)

	err := gopherLua.DoFile("example/luamodule/run.lua")
	if err != nil {
		c.JSON(http.StatusOK, formatError(err))
		return
	}
	err = gopherLua.ExecuteFunc("initParams", 1)
	if err != nil {
		c.JSON(http.StatusOK, formatError(err))
		return
	}
	ret := gopherLua.GetAndPop(-1) // returned value
	c.JSON(http.StatusOK, formatSuccess(goTool.TransLuaValue2Map(ret)))
}

//获取参数拓展
func setDefaultGinStatus(c *gin.Context, L *lua.LState) {
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
}
