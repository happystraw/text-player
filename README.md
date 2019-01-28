# Text Player

> A Small Tool for Converting Text to Speech and Playing Speech Based on Xunfei Web Api [讯飞在线语音合成(Xunfei online tts web api)](https://www.xfyun.cn/services/online_tts).

## Platforms

- Linux
- Windows
- macOS(Untested, I  don't have a macOS system computer. you can install from source)

## Installing

### Download binary

[Download from release](https://github.com/happystraw/text-player/releases)

### Install from source

```bash
go get -u github.com/happystraw/text-player
```

## Getting Started

### Configuration

- [Register a Xunfei online tts web api application](https://console.xfyun.cn/app/myapp). If you already have one, ignore this step.

- Add your host's IP to the application's IP whitelist.

- Set the `APPID` and `APIKey` of Xunfei online tts web api to default configuration for Text Player.

  ```bash
  ./text-player config --appid 'your xunfei api appid' --apikey 'your xunfei api apikey'
  ```

### Run in command line

```bash
./text-player play "小明同学, 早上好"
```

### Run in server

```bash
# default listen at localhost:8888
./text-player serve
```

**Test**

```bash
# POST  HTTP/1.1
# Host: localhost:8888
#
# Content-Disposition: form-data; name="msg"
#
# Hello, China
# ------WebKitFormBoundary7MA4YWxkTrZu0gW--
curl -d "msg=Hello, China" localhost:8888
```

## Usages

- Use "text-player --help" for more information about this application.

  <details>
    <summary>See detail</summary>

    <pre>
   ______        __  ___  __
  /_  __/____ __/ /_/ _ \/ /__ ___ _____ ____
   / / / -_) \ / __/ ___/ / _ `/ // / -_) __/
  /_/  \__/_\_\\__/_/  /_/\_,_/\_, /\__/_/
                              /___/
  Text Player is a Small Tool for Converting Text to Speech and Playing Speech.

  Usage:
    text-player [command]

  Available Commands:
    clean       Erase generated files
    config      Save global flags configuration to file
    help        Help about any command
    play        Text to speech and play it
    serve       Run on server

  Flags:
        --apikey string       xunfei tts api auth apikey
        --appid string        xunfei tts api auth appid
    -o, --cache-path string   path for cache files(default is $HOME/.text-player)
    -c, --config string       config file (default is $HOME/.text-player.yaml)
    -n, --disable-cache       disable cache generated speech files
    -h, --help                help for text-player
        --version             version for text-player

  Use "text-player [command] --help" for more information about a command.
    </pre>
  </details>

- Use "text-player [command] --help" for more information about a command. e.g. `./text-player play --help`

  <details>
    <summary>See detail</summary>

    <pre>
  Text to speech and play it

  Usage:
    text-player play &lt;message&gt; [flags]

  Flags:
    -h, --help   help for play

  Global Flags:
        --apikey string       xunfei tts api auth apikey
        --appid string        xunfei tts api auth appid
    -o, --cache-path string   path for cache files(default is $HOME/.text-player)
    -c, --config string       config file (default is $HOME/.text-player.yaml)
    -n, --disable-cache       disable cache generated speech files
    </pre>
  </details>

## License

Text Player is released under the Apache 2.0 license. See [LICENSE](./LICENSE)