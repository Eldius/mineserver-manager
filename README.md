
# mineserver-manager #

## Latest build status ##

![Go](https://github.com/eldius/mineserver-manager/actions/workflows/go.yml/badge.svg)
![Latest release status](https://github.com/eldius/mineserver-manager/actions/workflows/releases.yml/badge.svg)

A tool to help Minecraft servers management.

## reference links ##

- [APIs Mojang - Minecraft Wiki](https://minecraft.wiki/w/Mojang_API)
- [Mojang API - Wiki VG](https://wiki.vg/Mojang_API)


## snippets ##

```shell
## install command snippet 

mineserver \
    install \
      --version 1.21.3 \
      --motd "A new test server" \
      --level-name "My test world" \
      --enable-rcon \
      --headless \
      --dest . \
      --seed 5516949179205280665 \
      --memory-limit 2g \
      --whitelist-user Eldius
```

```http request
## get latest version files
GET https://api.github.com/repos/eldius/mineserver-manager/releases/latest
```

```shell
## download latest version by arch
wget $( curl https://api.github.com/repos/eldius/mineserver-manager/releases/latest | jq -r '.assets[].browser_download_url' | grep <my-arch> )
```
