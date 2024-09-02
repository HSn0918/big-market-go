-- script.lua

-- 请求的 URL
local url = "/api/v1/raffle/random_raffle"

-- 请求的 JSON 数据
local json_data = '{"strategy_id":100001}'

-- 设定请求的表头
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = json_data

-- 生成请求
request = function()
    return wrk.format(nil, url)
end
