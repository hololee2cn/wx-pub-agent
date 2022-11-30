# Wx-pub-agent

[![Code Coverage](https://img.shields.io/codecov/c/github/hololee2cn/wx-pub-agent/master.svg?style=flat-square)](https://app.codecov.io/gh/hololee2cn/wx-pub-agent/)
[![DOcker Pulls](https://img.shields.io/docker/pulls/leeoj2/pubplatform.svg)][hub]

中文 | [English](https://github.com/hololee2cn/wxpub/blob/master/doc/README_EN.md)

[hub]: https://hub.docker.com/repository/docker/leeoj2/pubplatform

## Introduction
>A WeChat public platform back-end service, currently mainly positioned as proxy push messages.

- **Support mass template message push**;
- **Support keyword text reply**;
- **Support user message push status query**;
- **Support localized storage of user information**;

## Architecture overview
![architecture overview](./doc/img/architecture.svg)

## Install

There are various ways of installing Wx-pub-agent.

### Docker images

Docker images are available on  [Docker Hub](https://hub.docker.com/r/leeoj2/pubplatform).

You can also use docker compose to [launch](./docker/docker-compose.yaml).

```bash
docker-compose up -d
```

Wx-pub-agent will now be reachable at <http://localhost:80/>

### Building and run from source

To build Prometheus from source code, You need:

* Go [version 1.17 or greater](https://golang.org/doc/install).

```bash
go build ./src/main.go
./main webapi
./main captcha
```

> Note: [Environment](./src/webapi/config/dev_configs.toml) is default dev.</br>
> [MYSQL table script](./docker/initsql) is here.


## Contact and Feedback
- We recommend that you use [github issue](https://github.com/hololee2cn/wxpub/issues) as the preferred channel for issue feedback and requirement submission;

## Contributing
We welcome you to participate in open source projects in various ways, including but not limited to:
- Feedback on problems and bugs encountered in use => [github issue](https://github.com/hololee2cn/wxpub/issues)
- Submit code to make the service faster, more stable and better =>[github PR](https://github.com/hololee2cn/wxpub/pulls)

## TODO
- [ ] Support pluggable sms service and captcha service
- [ ] Custom menu creation and event response (optional)
- [ ] One-click script deployment of sms service and captcha service

## License
[MIT](https://github.com/hololee2cn/wxpub/blob/master/LICENSE)