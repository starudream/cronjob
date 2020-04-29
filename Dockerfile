FROM golang:1.14-alpine AS builder1

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy go build -o exe .

FROM starudream/alpine:latest AS builder2

COPY --from=builder1 /build/exe /exe

RUN apk add --no-cache upx && upx -9 -q /exe

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY config.json config.json
COPY --from=builder2 /exe /exe

CMD /exe
