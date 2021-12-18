FROM golang:alpine3.15 as builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY . .
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update \
            && apk upgrade \
            && apk add --no-cache libwebp-dev build-base \
            && rm -rf /var/cache/apk/*
RUN go build -o MiraiGo .

FROM alpine:3.15 as runner
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update \
            && apk upgrade \
            && apk add --no-cache libwebp-dev \
            && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/MiraiGo \
 /app/application.yaml \
 /app/device.json \
 ./

ENTRYPOINT ["./MiraiGo"]