# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.16.3 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .


RUN rm -rf /etc/localtime && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # 安装protoc
    apt-get update && apt-get install -y unzip && \
    # go install github.com/golang/protobuf/protoc-gen-go@v1.26.0 && \
    go install github.com/gogo/protobuf/protoc-gen-gofast@v1.3.2 && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.17.1/protoc-3.17.1-linux-x86_64.zip && \
    unzip protoc-3.17.1-linux-x86_64.zip -d protoc && rm protoc-3.17.1-linux-x86_64.zip && \
    cp -r protoc/bin/* /usr/bin && cp -r protoc/include/* /usr/include && rm -rf protoc && \
    # 指定OS等，并go build
    GOOS=linux GOARCH=amd64 && \
    make && \
    # 由于我不止依赖二进制文件，还依赖config文件夹下的文件
    # 所以我将这些文件放到了publish文件夹
    mkdir publish && cp -r cmd publish && \
        cp -r config publish

# 运行阶段指定alpine作为基础镜像
FROM alpine

WORKDIR /app

ARG env=prod

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/publish .

RUN mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# 指定运行时环境变量
ENV ENV=${env}

EXPOSE 80 52898 6060

ENTRYPOINT ["./cmd/payment"]

