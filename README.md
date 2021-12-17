# rpcx-service
一个基于RPCX（Go语言的快、易用却功能强大RPC服务治理框架）搭建的RPC服务框架。 

特性：简单易用、超快高效、功能强大、服务发现、服务治理、服务分层、版本控制、路由标签注册。

## 框架目录介绍
- bin `脚本目录`
- log `日志目录`
- src `源码目录`
    - config `框架配置目录`
    - libraries `资源库目录`
    - logic `逻辑层目录`
    - model `数据层目录`
    - service `服务层目录`
    - table `数据表层目录`
    - go.mod `包版本定义`
    - main.go `程序入口`
    - service.json `服务配置文件`

## 服务如何跑起来？
1. 执行 `./bin/produce.sh help` 命令，查看脚手架帮助文档，执行对应子命令，会有相应信息输出。
    - 1.1. 执行 `./bin/produce.sh sync` 命令，同步服务vendor资源。
    - 1.2. 执行 `./bin/produce.sh build local` 命令，可以编译当前代码生成服务程序。
    - 1.3. 执行 `./bin/produce.sh reload` 命令，实现平滑重启服务。
    - 1.4. 执行 `./bin/produce.sh status` 命令，可以查看服务状态。
    - 1.5. 执行 `./bin/produce.sh quit` 命令，可以平滑停止当前服务。

## 参考文献
- RPCX文档：https://doc.rpcx.io
- Redis命令文档：http://redis.cn/commands.html
- GORM文档：https://gorm.io/zh_CN/docs
- Docker文档：http://www.dockerinfo.net/document


