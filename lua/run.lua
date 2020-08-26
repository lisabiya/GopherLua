function initParams()
    local example = require('db_module.example')
    local name = getParams("name")
    print(name)
    RouterF = {
        insertCol = example.insertCol,
        remove = example.remove,
        update = example.update,
        getList = example.getList,
    }
    if RouterF[name] then
        local code, tables = RouterF[name]()
        return { code = code, response = tables }
    else
        return { code = -1, response = "没有发现此函数" }
    end
end

