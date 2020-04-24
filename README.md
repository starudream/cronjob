# CronJob

![Go](https://github.com/starudream/cronjob/workflows/Go/badge.svg)
![Docker](https://github.com/starudream/cronjob/workflows/Docker/badge.svg)
![License](https://img.shields.io/badge/license-Apache%20License%202.0-blue)

## Config

```json
{
  "tasks": [
    {
      "name": "test",
      "url": "https://api.github.com/",
      "body": "",
      "cron": "* * * * *",
      "timezone": "Asia/Shanghai",
      "method": "GET",
      "headers": {},
      "timeout": 30
    }
  ]
}
```

## Usage

![Build](https://img.shields.io/docker/cloud/build/starudream/cronjob)
![Version](https://img.shields.io/docker/v/starudream/cronjob)
![Size](https://img.shields.io/docker/image-size/starudream/cronjob/latest)
![Pull](https://img.shields.io/docker/pulls/starudream/cronjob)

```bash
docker pull starudream/cronjob
```

```bash
docker run -d \
    --name cronjob \
    --restart always \
    -e DEBUG=true \
    -v /opt/docker/cronjob/config.json:/config.json \
    starudream/cronjob:latest
```

## License

[Apache License 2.0](./LICENSE)
