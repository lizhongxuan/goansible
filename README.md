# goansible
嵌入式ansible
- 服务端根据playbook配置跟主机及中间件参数进行模板拼接
- 分发给work去执行shell

![goansible.png](docs%2Fgoansible.png)

## 代码模块
- ansible: 负责解析playbook文件,按配置工作流分发任务给work
- module:  负责多个module的赋能,每个module支持一种命令格式
- work:    负责执行下发的shell任务

## example
- playbook.yaml: 记录运行的脚本内容跟逻辑关系
- host.yaml: 记录主机资源
- toolcheck.yaml: 记录工具的检查命令跟回调命令

## checker
- 通过配置编排检测循序
- 运行某个功能之前,通过配置,检查工具是否能使用
- 支持对输出进行正则匹配检查,校验不通过报错
- 支持编写hook,根据成功或者失败来运行回调脚本

## TODO
- 完善module基础模块
- 抽象module
- 实现多种work
- 引入泛型
- 优化执行work