# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.17-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 复制项目中的所有内容
COPY . .

# 构建 API 可执行文件，把文件输出到 /app/build
RUN mkdir -p /app/build
RUN go build -o /app/build/api ./backend/api/main.go

# 构建 Worker 可执行文件，把文件输出到 /app/build
RUN go build -o /app/build/worker ./backend/cmd/worker/main.go

# 使用一个更小的运行时镜像
FROM alpine:latest
WORKDIR /root/

# 从构建阶段中复制 API 和 Worker 二进制文件到最终镜像
COPY --from=builder /app/build/api .
COPY --from=builder /app/build/worker .

# 默认启动命令
CMD ["sh"]
