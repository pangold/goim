FROM ubuntu:latest

MAINTAINER pangold<pangold@163.com>

RUN mkdir /im/
WORKDIR /im/
COPY build/goim /im/

# GRPC, HTTP, JSONRPC
EXPOSE 9527 9528 9529

ENTRYPOINT ["./goim"]