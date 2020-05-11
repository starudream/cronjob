FROM golang:1.14-alpine AS builder1

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.cn,direct go build -o cronjob .

FROM starudream/alpine:latest AS builder2

COPY --from=builder1 /build/cronjob /cronjob

RUN apk add --no-cache upx && upx -9 -q /cronjob

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY config.json config.json
COPY --from=builder2 /cronjob /cronjob

CMD /cronjob
