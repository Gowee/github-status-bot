[![Badge for Docker build status](https://img.shields.io/docker/cloud/build/gowe/ghstsbot)](https://hub.docker.com/repository/docker/gowe/ghstsbot)
![Badge for CI Check workflow](https://github.com/Gowee/github-status-bot/workflows/CI%20Check/badge.svg)
# GitHub Status Bot
The [@ghstsbot](https://t.me/ghstsbot) that powers the [@GitHub_Status](https://t.me/GitHub_Status) channel on Telegram.
<!-- potential-octo-memory-->

**CLi**:
```shell
podman volume create ghstsbot
podman create --name ghstsbot -e TELEGRAM_BOT_TOKEN="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" -e CHAT_ID="-1001234567890" -e CHECK_INTERVAL=300 ghstsbot
```

----
<sub>The project, channel and bot is not affiliated in any way with GitHub. All images related to GitHub or Octocat are used for [fair use](https://en.wikipedia.org/wiki/Fair_use) only and not covered by the copyright license of the project.</sub>
