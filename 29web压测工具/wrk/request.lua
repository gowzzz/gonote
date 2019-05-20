wrk.method = "POST"
wrk.body   = "foo=bar&baz=quux"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

-- 为每次request更换一个参数
request = function()
    uid = math.random(1, 10000000)
    path = "/test?uid=" .. uid
    return wrk.format(nil, path)
end
-- 每次请求之前延迟10ms
function delay()
    return 10
end

-- 每个线程要先进行认证，认证之后获取token以进行压测
token = nil
path  = "/authenticate"

request = function()
   return wrk.format("GET", path)
end



response = function(status, headers, body)
   if not token and status == 200 then
      token = headers["X-Token"]
      path  = "/resource"
      wrk.headers["X-Token"] = token
   end
end

-- 压测支持HTTP pipeline的服务
-- 通过在init方法中将三个HTTP request请求拼接在一起，实现每次发送三个请求，以使用HTTP pipeline。
init = function(args)
    local r = {}
    r[1] = wrk.format(nil, "/?foo")
    r[2] = wrk.format(nil, "/?bar")
    r[3] = wrk.format(nil, "/?baz")
 
    req = table.concat(r)
 end
 
 request = function()
    return req
 end
