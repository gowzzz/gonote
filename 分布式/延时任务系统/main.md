
一般有两种思路来解决这个问题：
	实现一套类似crontab的分布式定时任务管理系统。
	实现一个支持定时发送消息的消息队列。

定时器的实现
	1.时间堆：最常见的时间堆一般用小顶堆实现，小顶堆其实就是一种特殊的二叉树
	2.时间轮
任务分发
数据再平衡和幂等考量	

定时任务管理系统
https://github.com/ouqiang/gocron
https://github.com/george518/PPGo_Job
https://www.cnblogs.com/jssyjam/p/11910851.html
有的是一个系统:

https://github.com/shunfei/cronsun
https://github.com/ouqiang/gocron
https://github.com/lisijie/webcron
有的是一个库:
https://github.com/robfig/cron

链接：https://www.jianshu.com/p/83f37db7b078


使用 etcd实现 分布式配置管理