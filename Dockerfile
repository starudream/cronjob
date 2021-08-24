FROM golang:alpine AS builder

WORKDIR /build

COPY . .

COPY --from=tonistiigi/xx:golang-master / /

RUN CGO_ENABLED=0 GO111MODULE=on go build -o cronjob .

FROM starudream/alpine-glibc:latest

COPY config.json config.json
COPY --from=builder /build/cronjob /cronjob

WORKDIR /

CMD /cronjob
