FROM starudream/golang AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on go build -o cronjob . && upx cronjob

FROM starudream/alpine-glibc:latest

COPY config.json config.json
COPY --from=builder /build/cronjob /cronjob

WORKDIR /

CMD /cronjob
