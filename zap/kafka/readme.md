# 先安装java
# kafka依赖zookeeper
http://zookeeper.apache.org/releases.html
3.5.5
打开zookeeper-3.4.13\conf，把zoo_sample.cfg重命名成zoo.cfg
把dataDir的值改成 dataDir=F:/tools/apache-zookeeper-3.5.5-bin/data
2.6 添加如下系统变量：
ZOOKEEPER_HOME: F:\tools\apache-zookeeper-3.5.5-bin (zookeeper目录)
Path: 在现有的值后面添加 ";%ZOOKEEPER_HOME%\bin;"
运行Zookeeper: 打开cmd然后执行 zkserver

# kafka安装
 http://kafka.apache.org/downloads.html
 http://mirror.bit.edu.cn/apache/kafka/2.3.0/kafka_2.12-2.3.0.tgz
 从文本编辑器里打开 server.properties
 log.dirs=./logs

 输入并执行:  .\bin\windows\kafka-server-start.bat .\config\server.properties


 #  创建TOPICS
 cd F:\tools\kafka_2.12-2.3.0\bin\windows
 创建一个topic： 
 kafka-topics.bat --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test

 # 打开一个PRODUCER
 kafka-console-producer.bat --broker-list localhost:9092 --topic test

 # 打开一个CONSUMER
 kafka-console-consumer.bat --bootstrap-server localhost:9092 --topic test --from-beginning