<center><img src="./logo.png"></center>



<center>一个基于`MiraiGo`项目的QQ群<-->TG Chat Group的消息同步机器人.资源占用低(10M左右内存),性能高效.</center>

## 🥗 功能介绍

* 消息转发

* 回复提醒

  - Telegram端由Bot API限制，使用缓存以发送消息的用户名称实现另端回复消息时@发送人.
  - QQ端通过正则匹配获取消息中的QQ号字段实现。

* 消息卡片解析

* Telegram Stickers转发

  * 出于技术限制，动态贴纸无法转发预览。可以部署基于 `Cloudflare Workers` 的 TGS 预览服务，详情请参照 [TGS Preview](https://github.com/lx200916/TGS_Preview)

* 临时屏蔽命令

  * 使用 /on /off来开关转发，在Caption和文本消息前加入// 前缀来屏蔽转发。屏蔽指在TG端开关转发，关闭后TG消息将不同步到QQ端，QQ端仍可同步到TG.

  *  使用临时屏蔽命令后，机器人可以通过为群名称加上🔔/🔕 Emoji来标识此刻转发状态。需给予管理员权限。

>注意:
>
>* 消息同步中包括将QQ音乐分享卡片转发到Telegram的功能,当分享的音乐来自国内音乐软件时,需要国内IP才能获取到音乐文件并转发到Telegram.如果部署在国外IP的VPS上,请禁用此功能.
>* Telegram Bot需要关闭隐私模式.
>* 如在Docker之外部署,请保证安装`libwebp`软件包,本项目需要它来将`Webp`图像转换为QQ支持的`JPEG`图像.
>* 为确保账号不触发风控，可以尝试先在国内IP家宽环境下运行本程序，将生成的`device.json`打包入镜像，再上传到云服务器。
>* 如果需要转发`Telegram`端动态贴纸，可以参考 [TGS_Preview](https://github.com/lx200916/TGS_Preview) 项目。

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



## 🌮 进阶内容

### 🐳 Docker 支持

参照 [Dockerfile](./Dockerfile).

```bash
git clone https://github.com/lx200916/Forward_Go.git
cd Forward_Go
mkdir data logs
docker build -t forwordgo .
docker run -d -v $pwd/data:/app/data -v $pwd/logs:/app/logs forwordgo
```



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
