--------------
创建服务器私钥，命令会让你输入一个口令：
openssl genrsa -des3 -out server.key 1024

创建签名请求的证书（CSR）：
user@ubuntu:bug28$ openssl req -new -key server.key -out server.csr
Enter pass phrase for server.key:
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:ch
State or Province Name (full name) [Some-State]:bj
Locality Name (eg, city) []:bj
Organization Name (eg, company) [Internet Widgits Pty Ltd]:bj
Organizational Unit Name (eg, section) []:zn
Common Name (e.g. server FQDN or YOUR name) []:bj
Email Address []:bj

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:1qaz2WSX
An optional company name []:bj


在加载SSL支持的Nginx并使用上述私钥时除去必须的口令：
# cp server.key server.key.org
# openssl rsa -in server.key.org -out server.key
最后标记证书使用上述私钥和CSR：
# openssl x509 -req -days 3650 -in server.csr -signkey server.key -out server.crt

/home/user/bug28

server{
		listen 80 default_server;
		server_name   10.58.122.238;
		return 301 https://$server_name$request_uri;
}
server {
		server_name   10.58.122.238;
		listen 443;
		ssl on;
		ssl_certificate /home/user/bug28/server.crt;
		ssl_certificate_key /home/user/bug28/server.key;

		charset utf-8;
		proxy_read_timeout 86400;
		error_log    /home/user/logs/nginx-error.log    error;



#               listen 443;
#               ssl on;
 #              ssl_certificate /home/user/bug28/server.crt;
  #             ssl_certificate_key /home/user/bug28/server.key;
#               ssl_session_timeout  5m;
       # ssl_protocols  SSLv2 SSLv3 TLSv1;
       # ssl_ciphers  HIGH:!aNULL:!MD5;
       # ssl_prefer_server_ciphers  on;

       # proxy_redirect http:// $scheme://;
       # port_in_redirect on;


-----------------------------
1.luajit lua的增强版本
wget http://luajit.org/download/LuaJIT-2.0.5.tar.gz
https://github.com/openresty/luajit2/releases //这是nginx专用的优化版本
修改makefile 
export PREFIX= /usr/local/luajit
make && make install

2.nginx安装
安装gcc g++的依赖库

sudo apt-get install build-essential
sudo apt-get install libtool
安装pcre依赖库（http://www.pcre.org/）

sudo apt-get update
sudo apt-get install libpcre3 libpcre3-dev
安装zlib依赖库（http://www.zlib.net）

sudo apt-get install zlib1g-dev
安装SSL依赖库（16.04默认已经安装了）

sudo apt-get install openssl




-------------------
listen 443 ssl http2;
server_name  localhost;
ssl_certificate_key /home/anfang/Downloads/cert/key.pem;
ssl_certificate   /home/anfang/Downloads/cert/cert.pem;

openssl genrsa -out key.pem 2048
openssl req -new -x509 -sha256 -key key.pem 
-out cert.pem -days 36500 
-subj /C=CN/ST=Shanghai/L=Songjiang/O=ztgame/OU=tech/CN=mydomain.ztgame.com/emailAddress=myname@ztgame.com
openssl x509 -in cert.pem -noout -text
----------------------
2）下载 ngx_devel_kit (NDK) 模块 . 解压

（3）下载 ngx_lua .解压
https://github.com/simplresty/ngx_devel_kit/tags
https://github.com/openresty/lua-nginx-module/tags


./configure --prefix=/opt/nginx-build-1.8.1 \
--with-ld-opt="-Wl,-rpath,/opt/luajit-2.0.0/lib" \
--add-module=/opt/lua-nginx-module-0.10.2 \
--add-module=/opt/ngx_devel_kit-0.3.0rc1

wget https://github.com/openssl/openssl/archive/OpenSSL_1_1_1a.tar.gz
tar zxvf OpenSSL_1_1_1a.tar.gz

cd nginx-1.15.8/
./configure --prefix=/usr/local/nginx --with-pcre --with-luajit   --with-http_stub_status_module --with-http_ssl_module --with-http_realip_module --with-http_v2_module --with-openssl=../openssl-OpenSSL_1_1_1a 	--with-ld-opt="-Wl,-rpath,/usr/local/luajit/lib"   --add-module=/usr/local/servers/lua-nginx-module-0.10.14 --add-module=/usr/local/servers/ngx_devel_kit-0.3.1rc1 -j2 



make -j2
sudo make install


---------------------
#lua指令方式
#在server 中添加一个localtion
location /hello {
            default_type 'text/plain';
            content_by_lua 'ngx.say("hello, lua")';
        }
#lua文件方式

#在server 中添加一个localtion
location /lua {
    default_type 'text/html';
    content_by_lua_file conf/lua/test.lua; #相对于nginx安装目录
}
#test.lua文件内容
ngx.say("hello world");

------------------------------------
在http部分添加如下配置 
#lua模块路径，多个之间”;”分隔，其中”;;”表示默认搜索路径，默认到/usr/servers/nginx下找  
lua_package_path "/usr/servers/lualib/?.lua;;";  #lua 模块  
lua_package_cpath "/usr/servers/lualib/?.so;;";  #c模块   

项目模块化
example
    example.conf     ---该项目的nginx 配置文件
    lua              ---我们自己的lua代码
      test.lua
    lualib            ---lua依赖库/第三方依赖
      *.lua
      *.so
=》
http {  
    include       mime.types;  
    default_type  text/html;  
  
    #lua模块路径，其中”;;”表示默认搜索路径，默认到/usr/servers/nginx下找  
    lua_package_path "/usr/example/lualib/?.lua;;";  #lua 模块  
    lua_package_cpath "/usr/example/lualib/?.so;;";  #c模块  
    include /usr/example/example.conf;  
-------example.conf-------
server {  
    listen       80;  
    server_name  _;  
  
    location /lua {  
        default_type 'text/html';  
        lua_code_cache off;  
        content_by_lua_file /usr/example/lua/test.lua;  
    }  
}  
