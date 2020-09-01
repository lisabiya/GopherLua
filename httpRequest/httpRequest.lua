---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wakfu.
--- DateTime: 2020/8/1 11:37
---
Http = {}

function Http.new()
    local obj = httpRequest.new()
    return setmetatable({ http = obj }, { __index = Http })
end


--*********************go实现的接口api**************************

--Post
---@param   url string   路径
---@param   params table form表单
---@return  number|table
function Http:postForm(url, params, type)
    return self.http:postForm({ url = url, params = params, type = type })
end

--Get
---@param   url string  路径
---@param   query table query参数
---@return  number|table
function Http:getQuery(url, query)
    return self.http:getQuery({ url = url, query = query })
end

return Http


