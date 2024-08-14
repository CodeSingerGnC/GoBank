# build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add --no-cache curl tar
RUN curl -L "https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz" | tar xvz

# run stage 
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main ./main
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration
EXPOSE 8080
ENV DB_CONNECTION_SOURCE="mysql://root:MySQLPassword@tcp(mysql8:3306)/microbank"

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

