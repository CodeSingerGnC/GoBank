# Description: Makefile for MicroBank
# Author: CodeSingerGnC
# Created: 2024-04-1
# Last update: 2024-08-5
# 本文件主要用于项目的开发、测试等操作，不适用于生产环境。

DB_URL = mysql://root:MySQLPassword@tcp(localhost:3306)/microbank

# container
createdk:
	docker run --name microbank -e MYSQL_ROOT_PASSWORD=MySQLPassword -p 3306:3306 -d mysql:8
stopdk:
	docker stop microbank
removedk:
	docker rm microbank
startdk:
	docker start microbank

# for docker compose test
cpup:
	docker-compose up
cpdown:
	docker-compose down
dkrm:
	docker-compose rm
rmimg:
	docker rmi microbank-api 

# database
createdb:
	docker exec -it microbank mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS microbank;"
dropdb:
	docker exec -it microbank mysql -u root -p -e "DROP DATABASE IF EXISTS microbank;"

# migrate
mgup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
mgdown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
mgdown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

# sqlc
sqlc:
	sqlc generate
vet:
	sqlc vet

# dbdocs
dbdocs:
	dbdocs build doc/dbdocs/db.dbml
schema:
	dbml2sql doc/dbdocs/db.dbml --mysql -o doc/dbdocs/schema.sql

# test
test:
	go test -v -cover -short -count=1 ./... || (echo "测试失败，请检查错误信息。" && false)

# main
run:
	go run main.go

# redis
redis:
	docker run --name microbank-redis -p 6379:6379 -d redis:alpine
redis-ping:
	docker exec -it microbank-redis redis-cli ping
redis-stop:
	docker stop microbank-redis
redis-rm:
	docker rm microbank-redis
redis-restart:
	docker restart microbank-redis

# mock
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/CodeSingerGnC/MicroBank/db/sqlc Store

# grpc 
protoc:
	rm -rf doc/swagger/*.swagger.json
	rm -rf pb/*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=microbank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans -r repl --host localhost --port 5403
	
.PHONY: createdk stopdk removedk startdk createdb dropdb mgup mgdown gen vet test run protoc redis dbdocs schema redis-ping redis-stop redis-rm redis-restart mock
