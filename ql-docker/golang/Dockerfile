FROM ubuntu:20.04
ENV TIMEZONE Asia/Shanghai
ENV PATH /usr/local/go/bin:/usr/local/protoc/bin:/root/go/bin:$PATH
COPY go1.17.5.tar.gz /usr/local/src/go1.17.5.tar.gz
COPY protoc-3.17.3-linux-x86_64.zip /usr/local/src/protoc-3.17.3-linux-x86_64.zip

RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list  \
    && ln -snf /usr/share/zoneinfo/$TIMEZONE /etc/localtime && echo $TIMEZONE > /etc/timezone \
	&& apt update && apt install -y locales wget g++ gcc libc6-dev make pkg-config unzip sudo vim \
	&& tar -C /usr/local/ -xzf /usr/local/src/go1.17.5.tar.gz \
	&& rm -rf /usr/local/src/go1.17.5.tar.gz \
	&& unzip /usr/local/src/protoc-3.17.3-linux-x86_64.zip -d /usr/local/protoc/ \
	&& rm -rf /usr/local/src/protoc-3.17.3-linux-x86_64.zip \
	&& go env -w GOPROXY=https://goproxy.cn,direct \
	&& go env -w GO111MODULE=on
#	&& go install github.com/golang/protobuf/protoc-gen-go@latest
WORKDIR /www
