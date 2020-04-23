FROM golang:1.14-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.io,direct go build -o exe . \
    && sed -i 's|http://dl-cdn.alpinelinux.org|https://mirrors.tuna.tsinghua.edu.cn|g' /etc/apk/repositories \
    && apk add --no-cache upx \
    && upx -9 -q exe

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY config.json config.json
COPY --from=builder /build/exe exe

CMD /exe
