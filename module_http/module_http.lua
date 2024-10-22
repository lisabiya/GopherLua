---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wakfu.
--- DateTime: 2020/8/1 11:37
---
Http = {}

function Http.new()
    local obj = httpRequest.new()
    return setmetatable({ http = obj, options = {} }, { __index = Http })
end


--*********************go实现的接口api工具类，便于提供代码提示**************************

---@param targetUrl string
function Http:get(targetUrl)
    self.options["get"] = targetUrl
    return self
end

---@param targetUrl string
function Http:post(targetUrl)
    self.options["post"] = targetUrl
    return self
end

--  header
---@param content table  {string=string}
function Http:set(content)
    self.options["set"] = content
    return self
end

---@param content table
function Http:query(content)
    self.options["query"] = content
    return self
end

--    "text/html" uses "html"
--    "application/json" uses "json"
--    "application/xml" uses "xml"
--    "text/plain" uses "text"
--    "application/x-www-form-urlencoded" uses "urlencoded", "form" or "form-data"
---@param content string
function Http:type(content)
    self.options["type"] = content
    return self
end

--  form表单，json等
---@param content table
function Http:send(content)
    self.options["send"] = content
    return self
end

-- 启动查询
---@return  number|table  返回状态值和返回值
function Http:End()
    return self.http:End(self.options)
end

return Http


