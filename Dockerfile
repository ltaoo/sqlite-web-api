# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20

RUN apt-get update && apt-get install -y gcc libc6-dev

# 设置工作目录
WORKDIR /app

# 将当前目录的内容复制到容器内的 /app 目录
COPY . .

# 设置环境变量
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

ENV GIN_MODE=release

# 编译二进制文件
RUN go build -ldflags '-extldflags "-static" -X main.env=prod' -o sqliteweb

