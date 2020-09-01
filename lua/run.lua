function initParams()
    local example = require('db_module.example')
    local httpExample = require('httpRequest.httpExample')

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