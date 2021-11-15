FROM starudream/golang AS builder

WORKDIR /build

COPY . .

RUN apk add --no-cache make git \
    && make build \
    && make upx

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY config.json config.json

COPY --from=builder /build/bin/app /app

CMD /app
