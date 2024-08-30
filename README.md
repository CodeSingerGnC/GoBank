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

本项目主要涉及 5 张表： users, accounts, entries, transfers, session，其具体结构可以通过在线文档 [MicroBankDB](https://dbdocs.io/outof2023/MicroBankDB) 查看（此文档由 [dbdocs](https://dbdocs.io) 代理部署，第一次访问存在延迟，请耐心等待）。

## 项目重要部分详细介绍

[基于 bcrypt 算法的密码存储与验证](./doc/intro/password.md)

[基于 access_token 和 refresh_token 的鉴权](./doc/intro/Authtication.md)

[使用异步任务队列 asynq 发送邮件](./doc/intro/asynq.md)

[参数验证 API](./doc/intro/validator.md)

## 运行项目

### 运行服务需要使用的端口

- MySQL：3306
- gRPC: 5403
- HTTP: 8080

请在运行此项目之前确保端口处于可用状态。

### 通过 Docker Compose 运行服务

确保你已经安装了 Docker 和 Docker Compose。如果没有，请参考 [Docker 官方文档](https://docs.docker.com/get-docker/) 进行安装。

克隆仓库：

```sh
git clone git@github.com:CodeSingerGnC/GoBank.git
```

切换到 GoBank 工作目录

```sh
cd GoBank
```

使用 docker-compose 启动服务：

```sh
docker-compose up
```

![演示](doc/vhs/demonstration.gif)