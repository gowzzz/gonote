#!/usr/local/openresty/bin/resty
local mysql = require "resty.mysql"

local db,err=mysql:new()
if not db then
    ngx.say("new mysql failed:",err)
    return
end
db:set_timeout(1000)
-- local ok,err=db:set_keepalive(1000,10)
-- local count,err=db:get_reused_times()
local opts={
    host="127.0.0.1",
    port=3306,
    database="test",
    user="root",
    password="123456",
}
local ok,err=db:connect(opts)
if not ok then
    ngx.say("connect mysql failed:",err)
    db:close()
    return
end

ngx.say("version:",db:server_ver())

-- res,err=db:query(stmt,nrows) nrows是预估的结果行，res保存了查询结果和nil
-- local ok,err=db:query("create table test2(id int,name varchar(5))")
-- if not ok then
--     ngx.say("create test2 failed:",err)
--     db:close()
--     return
-- end
-- local res,err=db:query("insert into user values(3,'wocao')")
-- if err~=nil then
--     ngx.say("insert user failed:",err)
--     return
-- -- end
-- ngx.say(res.insert_id)
-- ngx.say(res.affected_rows)
-- local s = ngx.quote_sql_str("id,name")
local s = "id,name"
local selectStr="select "..s.." from user;"
ngx.say("select str:",selectStr)
local results,err=db:query(selectStr,5)
if err~=nil then
    ngx.say("select user failed:",err)
    return
end
for i,rows in ipairs(results) do 
    for k,v in pairs(rows) do 
        ngx.print(k,"=",v,";\n")
    end
end







db:close()

