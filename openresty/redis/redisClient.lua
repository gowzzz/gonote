#!/usr/local/openresty/bin/resty
local redis=require "resty.redis"
local rds,err=redis:new()
rds:set_timeout(1000) --1000ms
ok,err=rds:connect("127.0.0.1",3306)
if not ok then 
    ngx.say("fialed to connect :",err)
    rds:close()
    return 
end 

ok,err=rds:set_keepalive(1000,10) --放入连接池
count,err=rds:getreused_times()--获取复用的次数
-- res,err=rds:command(key,...)--如果命令执行结果是null，会返回ngx.null常量，即执行成功返回的是空结果
res,err=rds:get("a")
assert(res~=ngx.null)
ngx.say(res)

-- ok,err=rds:hset("zelda","bow",2017)
-- res,err=rds:hget("zelda","bow")

-- ok,err=rds:del("list")
-- ok,err=rds:lpush("list",1,2,3,4)
-- res,err=rds:lpop("list")
-- ngx.say("list:",res)

--当有大量操作时，可以使用管道  这是典型的空间换时间
-- rds:init_pipeline(n)--启动管道 n是预估的命令的数量
--正常发送命令
-- results,err=rds:commit_pipeline()--提交管道
-- rds:cancel_pipeline()--取消管道
-- rds:init_pipeline(10)
-- for i=1,10 do 
--     ok,err=rds:rpush("numbers",i)
--     if not ok then
--         rds:cancel_pipeline()
--     end 
-- end 
-- results,err=rds:commit_pipeline()

--事务 multi/exec 

--脚本 redis2.6以上内嵌支持lua，可以把复杂的事务过程编写成lua脚本，让业务逻辑泡在redis内部，消除通信交互的成本，
