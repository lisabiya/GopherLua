---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by wakfu.
--- DateTime: 2020/8/4 17:58
---
local example = {}
--
local json = require('module_tools.json')
local Http = require('module_http.module_http')

function example.get()
    local request = Http.new()
    local code, response = request:get("https://www.wanandroid.com/hotkey/json")
                                  :send({ name = "sss" })
                                  :query({ name = "nihao" })
                                  :End()
    if code == 0 then
        local ok, result = pcall(json.decode, response)
        return ok and 0 or -1, result
    else
        return code, response
    end
end

function example.post()
    local request = Http.new()
    local code, response = request:post("http://0.0.0.0:8080/api/invoice/order")
                                  :type("form")
                                  :send({ platforms = 2, account = "13000000000" })
                                  :set({ Session = "111" })
                                  :End()
    if code == 0 then
        local ok, result = pcall(json.decode, response)
        return ok and 0 or -1, result
    else
        return code, response
    end
end

return example