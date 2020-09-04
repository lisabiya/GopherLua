package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lisabiya/GopherLua"
	"github.com/lisabiya/GopherLua/example/controller"
	"github.com/lisabiya/GopherLua/module_http"
)

func main() {
	initRouter()
	//httpSimpleTest()
}

func httpSimpleTest() {
	gopherLua := GopherLua.NewState()
	gopherLua.Register(module_http.ModuleHttp{})

	err := gopherLua.DoString(`
  --直接引用声明模块
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

func initRouter() {
	r := gin.Default()
	controller.LoadLuaState()
	r.GET("/ping", controller.LoadLuaModule)
	_ = r.Run() // listen and serve on 0.0.0.0:8080
}
