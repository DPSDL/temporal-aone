# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.17-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制项目中的所有内容
COPY . .

# 构建 API 可执行文件
RUN go build -o /app/api ./main.go

# 构建 Worker 可执行文件
RUN go build -o /app/worker ./main.go

# 使用一个更小的运行时镜像
FROM alpine:latest

WORKDIR /root/

# 从构建步骤中复制 API 和 Worker 二进制文件
COPY --from=builder /app/api .
COPY --from=builder /app/worker .

# 默认启动命令，这里我们不指定，因为在 docker-compose.yml 中会指定具体的服务启动
CMD ["sh"]
