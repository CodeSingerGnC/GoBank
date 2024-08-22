# GoBank

此项目是一个基于 gRPC、MySQL、Docker、Github CI/CD 搭建的 GoLang 银行服务器。

## 技术栈

> 本项目开发使用到的所有框架与工具

- 开发语言：GoLang
- RPC 框架：grpc
- HTTP 代理：grpc-gateway
- gRPC 交互客户端: evans
- 基于 Redis 的异步任务队列：asynq
- 容器：Docker
- 简化常用开发命令：Makefile
- 版本控制：git
- 配置框架：viper
- 通过 sql 生成 GoLang 代码: sqlc
- token: JWT, PASETO
- 测试框架: testify, mock
- API 测试工具：postman
- CI/CD: GitHub Actions
- API 文档: swagger

## 项目表库结构

本项目主要涉及 5 张表： users, accounts, entries, transfers, session，其具体结构可以通过下面链接查看。

[MicroBankDB](https://dbdocs.io/outof2023/MicroBankDB)


## 运行项目


因为此服务需要在 3306 端口运行 MySQL，在 5403 端口运行 gRPC 服务，在 8080 端口运行 HTTP 服务，所以请在运行此项目之前确保端口处于可用状态。

并且需要确保你已经安装了 Docker 和 Docker Compose。如果没有，请参考 [Docker 官方文档](https://docs.docker.com/get-docker/) 进行安装。

1. 克隆仓库：
    ```sh
    git clone git@github.com:CodeSingerGnC/GoBank.git
    cd GoBank
    ```

2. 启动服务：

    ```sh
    docker-compose up
    ```

![演示](doc/vhs/demonstration.gif)

3. 访问服务：
    - gRPC 服务在 `localhost:5403`
    - HTTP 接口在 `localhost:8080`
    - MySQL 数据库在 `localhost:3306`
