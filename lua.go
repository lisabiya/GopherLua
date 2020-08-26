package main

import (
	"GopherLua/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	initRouter()
}

func initRouter() {
	r := gin.Default()
	r.GET("/ping", controller.LoadLuaModule)
	_ = r.Run() // listen and serve on 0.0.0.0:8080
}
