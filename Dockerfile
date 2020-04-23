FROM golang:1.14-alpine AS builder

WORKDIR /build

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o cronjob .

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY config.json config.json
COPY --from=builder /build/cronjob cronjob

CMD /cronjob
