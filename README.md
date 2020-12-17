# CronJob

![Go](https://img.shields.io/github/workflow/status/starudream/cronjob/Go/master?style=for-the-badge)
![Docker](https://img.shields.io/github/workflow/status/starudream/cronjob/Docker/master?style=for-the-badge)
![License](https://img.shields.io/badge/License-Apache%20License%202.0-blue?style=for-the-badge)

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
            "timeout": 30,
            "http_proxy": "http://127.0.0.1:7890",
            "https_proxy": "http://127.0.0.1:7890"
        }
    ]
}
```

## Usage

![Version](https://img.shields.io/docker/v/starudream/cronjob?style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/cronjob/latest?style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/cronjob?style=for-the-badge)

```bash
docker pull starudream/cronjob
```

```bash
docker run -d \
    --name cronjob \
    --restart always \
    -v /opt/docker/cronjob/config.json:/config.json \
    starudream/cronjob:latest
```

## License

[Apache License 2.0](./LICENSE)
