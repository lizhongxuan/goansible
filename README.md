# goansible
嵌入式ansible
- 服务端根据playbook配置跟主机及中间件参数进行模板拼接
- 分发给work去执行shell

## 代码模块
- ansible: 负责解析playbook文件,按配置工作流分发任务给work
- module:  负责多个module的赋能,每个module支持一种命令格式
- work:    负责执行下发的shell任务

## TODO
- 完善module基础模块
- 抽象module
- 实现多种work
- 引入泛型
- 优化执行work