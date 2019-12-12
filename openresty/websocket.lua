local server=require "resty.websocket.server"
local wb,err=server.new{
	timeout=5000,  --超时5s
	max_payload_len=1024*64  --数据帧最大64k
}
if not wb then
	ngx.log(ngx.ERR,"failed to init:",err)
	ngx.exit(444)
end

local data,typ,bytes,err
while true do 
    if not data then
        -- 空消息, 发送心跳
        local bytes, err = wb:send_ping()
        if not bytes then
          ngx.log(ngx.ERR, "failed to send ping: ", err)
          return ngx.exit(444)
        end
        ngx.log(ngx.ERR, "send ping: ", data)
    elseif typ == "close" then
        -- 关闭连接
        break
    elseif typ == "ping" then
        -- 回复心跳
        local bytes, err = wb:send_pong()
        if not bytes then
            ngx.log(ngx.ERR, "failed to send pong: ", err)
            return ngx.exit(444)
        end
    elseif typ == "pong" then
        -- 心跳回包
        ngx.log(ngx.ERR, "client ponged")
    elseif typ == "text" then
        -- 将消息发送
        bytes,err=wb:send_text(data)
    end
end