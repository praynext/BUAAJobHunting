FROM golang:latest AS builder
RUN mkdir /build
WORKDIR /build
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN swag init
RUN go build -o BUAAJobHunting .

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update
COPY --from=builder /build/BUAAJobHunting /app/BUAAJobHunting
COPY ./config.yaml /app/config.yaml
EXPOSE 9000
RUN mkdir ./log
CMD ./BUAAJobHunting
