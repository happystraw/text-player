# Text Player

> A Small Tool for Converting Text to Speech and Playing Speech Based on Xunfei Web Api [讯飞在线语音合成(Xunfei online tts web api)](https://www.xfyun.cn/services/online_tts).

## Platforms

- Linux
- Windows/macOS(Untested)

## Installing

```bash
go get -u github.com/happystraw/text-player
```

## Getting Started

### Configuration

- [Register a Xunfei online tts web api application](https://console.xfyun.cn/app/myapp). If you already have one, ignore this step.

- Set the `APPID` , `APIKey` and `APISecret` of Xunfei online tts web api to default configuration for Text Player.

  ```bash
  ./text-player config --appid 'your appid' --apikey 'your apikey' --apisecret 'your apisecret'
  ```

### Run in command line

```bash
./text-player play "富强、民主、文明、和谐、自由、平等、公正、法治、爱国、敬业、诚信、友善"
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


## License

Text Player is released under the Apache 2.0 license. See [LICENSE](./LICENSE)