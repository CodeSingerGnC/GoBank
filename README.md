# GoBank

此项目是一个基于 gRPC、MySQL、Docker、Github CI/CD 搭建的 GoLang 微服务架构银行服务器。

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
