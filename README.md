[![Badge for Docker build status](https://img.shields.io/docker/cloud/build/gowe/ghstsbot)](https://hub.docker.com/repository/docker/gowe/ghstsbot)
![Badge for CI Check workflow](https://github.com/Gowee/github-status-bot/workflows/CI%20Check/badge.svg)
# GitHub Status Bot
The [@ghstsbot](https://t.me/ghstsbot) that powers the [@GitHub_Status](https://t.me/GitHub_Status) channel on Telegram.
<!-- potential-octo-memory-->

**CLi**:
```shell
podman volume create ghstsbot
echo "📢📢📢
%s

Powered by https://git.io/ghstsbot" > /tmp/cdt.txt # optionally
podman create --name ghstsbot \
    -e TELEGRAM_BOT_TOKEN="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" \
    -e CHAT_ID="-1001234567890" \
    -e CHECK_INTERVAL=300 \
    -e CHAT_DESCRIPTION_TEMPLATE=$(cat /tmp/cdt.txt | base64 -w0)
    -v ghstsbot:/app/data \
    ghstsbot
podman start ghstsbot
podman logs -f ghstsbot
```
`CHECK_INTERVAL` defaults to be 300 if not set.

If `CHAT_DESCRIPTION_TEMPLATE` is not set, chat description updating will be disabled. Base64 encoding there is optional in case line break is not passed properly by shell. If present, `%s` should be included and will be replaced to generated content. Set it to a single `%s` (i.e. `-e CHAT_DESCRIPTION_TEMPLATE="%s"`) to activate chat description updating without putting any additional text.

----
<sub>The project, channel and bot is not affiliated in any way with GitHub. All images related to GitHub or Octocat are used for [fair use](https://en.wikipedia.org/wiki/Fair_use) only and not covered by the copyright license of the project.</sub>
