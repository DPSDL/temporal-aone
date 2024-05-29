# 使用官方 Go 镜像作为构建镜像
FROM golang:1.17 as build

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录中
COPY go.mod go.sum ./

# 安装依赖
RUN go mod download

# 将源代码复制到工作目录中
COPY . .

# 构建应用程序
RUN go build -o main .

# 使用一个更小的镜像作为运行时
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 将编译好的二进制文件从构建阶段复制到运行时镜像
COPY --from=build /app/main .

# 设置容器启动时运行的命令
CMD ["./main"]