## 引入最新的golan ，不设置版本即为最新版本
FROM golang as builder

MAINTAINER  <hao.wang>
## 在docker的根目录下创建相应的使用目录
## 设置工作目录
WORKDIR /go/hao
## 把当前（宿主机上）目录下的文件都复制到docker上刚创建的目录下
COPY . .
## 编译
RUN go env -w GO111MODULE=auto && go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod init gometrics && go mod tidy

RUN go build main.go

FROM ubuntu

WORKDIR /go/hao

COPY --from=builder /go/hao/main .

RUN apt-get update && apt-get -y  install sshpass
## 设置docker要开发的哪个端口
EXPOSE 18888
## 启动docker需要执行的文件
CMD ./main
