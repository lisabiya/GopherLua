function initParams()
    --当修改文件后需要移除已加载的包重新加载
    package.loaded["example/luamodule/exampleDB"] = nil
    --local tools = require('module_tools.tools')
    --tools.showLoadedPackage()
    --业务
    local example = require('example/luamodule/exampleDB')
    local httpExample = require('example.luamodule.exampleHttp')

    local name = getParams("name")
    RouterF = {
        insertCol = example.insertCol,
        remove = example.remove,
        update = example.update,
        getList = example.getList,
        get = httpExample.get,
        post = httpExample.post,
    }
    if RouterF[name] then
        local code, tables = RouterF[name]()
        return { code = code, response = tables }
    else
        return { code = -1, response = "没有发现此函数" }
    end
end