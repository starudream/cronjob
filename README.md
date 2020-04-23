# CronJob

[![GithubAction](https://github.com/starudream/cronjob/workflows/Go/badge.svg)](https://github.com/starudream/cronjob/actions)
[![License](https://img.shields.io/badge/license-Apache%20License%202.0-blue)](./LICENSE)

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

```bash
docker pull starudream/cronjob

docker run -d \
    --name cronjob \
    --restart always \
    -e DEBUG=true \
    -v `pwd`/config.json:/config.json \
    starudream/cronjob:latest
```

## License

[Apache License 2.0](./LICENSE)
