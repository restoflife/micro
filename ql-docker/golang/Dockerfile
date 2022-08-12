FROM ubuntu:latest

MAINTAINER dengquanguo

ENV TIMEZONE Asia/Shanghai
ENV PATH /usr/local/go/bin:/usr/local/protoc/bin:/root/go/bin:$PATH

COPY go1.19.linux-amd64.tar.gz /usr/local/src/go1.19.linux-amd64.tar.gz
COPY protoc-21.5-linux-x86_64.zip /usr/local/src/protoc-21.5-linux-x86_64.zip

RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && ln -snf /usr/share/zoneinfo/$TIMEZONE /etc/localtime && echo $TIMEZONE > /etc/timezone \
    && apt clean  \
    && apt update \
    && apt-get update -y \
    && apt install -y locales wget g++ gcc libc6-dev make pkg-config unzip sudo vim \
    && tar -C /usr/local/ -zxvf /usr/local/src/go1.19.linux-amd64.tar.gz \
	&& rm -rf /usr/local/src/go1.19.linux-amd64.tar.gz \
    && unzip /usr/local/src/protoc-21.5-linux-x86_64.zip -d /usr/local/protoc/ \
    && rm -rf /usr/local/src/protoc-21.5-linux-x86_64.zip \
	&& go env -w GOPROXY=https://goproxy.cn,direct \
	&& go env -w GO111MODULE=on \
    && go install github.com/golang/protobuf/protoc-gen-go@latest
WORKDIR /www