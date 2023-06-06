FROM golang:alpine

WORKDIR /app

# 将代码复制到容器中
COPY . .

# install make
RUN apk add --no-cache make && go mod download

# 编译
RUN go build -o ./bin/proxy

# 运行
CMD ["./bin/proxy", "server"]