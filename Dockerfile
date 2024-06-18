# 使用官方的 Golang 镜像作为基础镜像，指定所需Golang版本
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 复制项目中的所有内容
COPY . .

# 单独创建目录，以便看到具体错误
RUN mkdir -p ./build

# 列出当前目录内容
RUN ls -al /app

# 构建 API 可执行文件，把文件输出到 /app/build 目录，并输出详细调试信息
RUN go build -v -o ./build/api ./backend/api/main.go

# 构建 Worker 可执行文件，把文件输出到 /app/build 目录，并输出详细调试信息
RUN go build -v -o ./build/worker ./backend/cmd/worker/main.go

# 使用一个更小的运行时镜像
FROM alpine:latest

# 创建工作目录
WORKDIR /root/

# 从构建阶段中复制 API 和 Worker 二进制文件到最终镜像
COPY --from=builder /app/build/api .
COPY --from=builder /app/build/worker .

# 默认启动命令，这里我们不指定，因为在 docker-compose.yml 中会指定具体的服务启动
CMD ["sh"]
