# MiraiGo-TGForward

一个基于`MiraiGo`项目的QQ群<-->TG Chat Group的消息同步机器人.资源占用低(10M左右内存),性能高效.

## 🍱 基础配置

账号配置[application.yaml](./application.example.yaml)

```yaml
Telegram:
  #Bot Token
  token:
  #反代的Bot API 如https://xxx.xxx.com,不需要可留空
  APIAddr:
  #如果有Telegram贴纸预览服务,可在此写入,如 http://xx.xx.com,详情请参照 TGS_Preview 项目.
  TGSAddr:

bot:
  #账号
  account:
  #密码
  password:

Groups:
  #{QQ: QQ群号码, TG:Telegram Chat ID}
  - { QQ: , TG: }
```

> 注意:
>
> * 消息同步中包括将QQ音乐分享卡片转发到Telegram的功能,当分享的音乐来自国内音乐软件时,需要国内IP才能获取到音乐文件并转发到Telegram.如果部署在国外IP的VPS上,请禁用此功能.
> * Telegram Bot需要关闭隐私模式.
> * 如在Docker之外部署,请保证安装`libwebp`软件包,本项目需要它来将`Webp`图像转换为QQ支持的`JPEG`图像.
> * 为确保账号不触发风控，可以尝试先在国内IP家宽环境下运行本程序，将生成的`device.json`打包入镜像，再上传到云服务器。
> * 如果需要转发`Telegram`端动态贴纸，可以参考 [TGS_Preview](https://github.com/lx200916/TGS_Preview) 项目。

## 🌮 进阶内容

### 🐳 Docker 支持

参照 [Dockerfile](./Dockerfile)

## 🌏  引入的第三方 go module

- [MiraiGo](https://github.com/Mrs4s/MiraiGo)
  核心协议库

- [viper](https://github.com/spf13/viper)
  用于解析配置文件，同时可监听配置文件的修改

- [logrus](github.com/sirupsen/logrus)
  功能丰富的Logger

- [asciiart](github.com/yinghau76/go-ascii-art)
  用于在console显示图形验证码

- [telebot](https://github.com/tucnak/telebot)
  用于调用`Telegram Bot API`
