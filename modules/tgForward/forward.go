package tgForward

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/kolesa-team/go-webp/decoder"
	_ "github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	cmap "github.com/orcaman/concurrent-map"
	tb "gopkg.in/tucnak/telebot.v2"
	"image/png"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tgForward struct {
}

var logger = utils.GetModuleLogger("external.TGForward")

type QGroupInfo struct {
	QQNumber int64
	Flag     bool
}

var regQQReply = regexp.MustCompile(`\((\d+)\) :`)
var regQQReplyPhoto = regexp.MustCompile(`\((\d+)\)`)

//var regTGReply = regexp.MustCompile("-------\\s+(?U)(.+) :")

var GroupMap = make(map[int64]*tb.Chat)     //QQGroup-> TGChat Object
var GroupFlag = make(map[int64]*QGroupInfo) //TGGroup-> Flag

type groupConfig struct {
	QQ int64 `yaml:"QQ"`
	TG int64 `yaml:"TG"`
}

var MDReplace = strings.NewReplacer("_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(", "\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>", "#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|", "\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!")
var emojiRegExp = regexp.MustCompile(`[\x{1F515}\x{1F514}]$`)
var tgsAddr string

func (m *tgForward) Init() {
	var err error
	Bot, err = tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL:    config.GlobalConfig.GetString("Telegram.APIAddr"),
		Token:  config.GlobalConfig.GetString("Telegram.token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	tgsAddr = config.GlobalConfig.GetString("Telegram.TGSAddr")
	Bot.Handle("/on", func(m *tb.Message) {
		GroupFlag[m.Chat.ID].Flag = true
		_, _ = Bot.Send(m.Chat, "通讯已恢复")
		title := m.Chat.Title
		match := emojiRegExp.FindStringIndex(title)
		if match == nil {
			title = title + "🔔"
		} else {
			title = emojiRegExp.ReplaceAllString(title, "🔔")
		}
		err := Bot.SetGroupTitle(m.Chat, title)
		if err != nil {
			log.Println(err)
		}

	})
	Bot.Handle("/off", func(m *tb.Message) {
		GroupFlag[m.Chat.ID].Flag = false
		_, _ = Bot.Send(m.Chat, "不明力量截断了电波")
		title := m.Chat.Title
		match := emojiRegExp.FindStringIndex(title)
		if match == nil {
			title = title + "🔕"
		} else {
			title = emojiRegExp.ReplaceAllString(title, "🔕")
		}
		err := Bot.SetGroupTitle(m.Chat, title)
		if err != nil {
			log.Println(err)
		}

	})

	var grouplists []groupConfig
	err = config.GlobalConfig.UnmarshalKey("Groups", &grouplists)
	if err != nil {
		logger.Error("Group Lists Fail to Parse")
	}

	for _, groupInfo := range grouplists {
		tgChat, err := Bot.ChatByID(strconv.FormatInt(groupInfo.TG, 10))
		if err != nil {
			logger.Error("Group Fail to Parse:", groupInfo.TG)
		}

		GroupMap[groupInfo.QQ] = tgChat
		MembersMap[groupInfo.TG] = cmap.New()
		GroupFlag[groupInfo.TG] = &QGroupInfo{
			Flag: true, QQNumber: groupInfo.QQ,
		}

	}
	UpdateMembers()
	t := time.NewTicker(time.Duration(time.Minute * 5))

	//defer t.Stop()
	go func() {
		for {
			<-t.C
			SerializeMembers()
		}
	}()

}

func (m *tgForward) PostInit() {
	logger.Info("Init Finish")
}
func getReplyText(m *tb.Message) (string, error, *message.AtElement) {
	var at *message.AtElement = nil

	if strings.HasPrefix(m.Text, "//") || strings.HasPrefix(m.Caption, "//") {
		return "", errors.New("stop"), at
	}
	reply := m
	typeStr := ""
	if reply.Photo != nil {
		typeStr += "[Photo] "
	}
	if reply.Sticker != nil {
		typeStr += "[Sticker] "
	}
	if reply.Sender.IsBot {

		text := reply.Text
		if len(reply.Text) == 0 {
			text = reply.Caption

		}
		replyList := strings.Split(text, " -------- \n")
		result1 := regQQReply.FindStringSubmatch(replyList[len(replyList)-1])
		if len(result1) > 1 {
			atUid, _ := strconv.ParseInt(result1[1], 10, 64)
			at = message.NewAt(atUid)

		} else {
			if reply.Photo != nil {
				result1 = regQQReplyPhoto.FindStringSubmatch(replyList[len(replyList)-1])
				if result1 != nil {
					atUid, _ := strconv.ParseInt(result1[1], 10, 64)
					at = message.NewAt(atUid)

				}
			}
		}

	}
	replyText := fmt.Sprintf("%s %s :%s%s %s\n-------\n", reply.Sender.FirstName, reply.Sender.LastName, typeStr, reply.Text, reply.Caption)
	return replyText, nil, at
}
func midFWCancel(m *tb.Message) error {
	if m.Chat.Type != tb.ChatGroup && m.Chat.Type != tb.ChatSuperGroup {
		return errors.New("stop")
	}
	//QInfo, ok := GroupFlag[m.Chat.ID]
	//if !ok || (ok && !QInfo.Flag) {
	//	return ""
	//}
	if strings.HasPrefix(m.Text, "//") || strings.HasPrefix(m.Caption, "//") {
		return errors.New("stop")
	}
	return nil
}
func middleware(m *tb.Message, forwardText string) (string, error) {
	err := midFWCancel(m)
	UpdateAMember(m.Chat.ID, TGMember{UserID: m.Sender.ID, FirstName: m.Sender.FirstName, LastName: m.Sender.LastName, UserName: m.Sender.Username})
	if err != nil {
		return "", errors.New("stop")
	}
	forwardText = fwUserName(m, forwardText)
	forwardText = fwTGSource(m, forwardText)
	return forwardText, nil
}

func (m *tgForward) Serve(bot *bot.Bot) {
	//go getUpdates(Bot.Updates,bot)
	Bot.Handle(tb.OnNewGroupTitle, func(m *tb.Message) {
		if m.Sender.ID == Bot.Me.ID {
			err := Bot.Delete(m)
			if err != nil {
				log.Println(err)

			}
		}
	})
	Bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.Chat.Type != tb.ChatGroup && m.Chat.Type != tb.ChatSuperGroup {
			return
		}
		QInfo, ok := GroupFlag[m.Chat.ID]
		if !ok || (ok && !QInfo.Flag) {
			return
		}
		//if strings.HasPrefix(m.Text, "//") {
		//	return
		//}
		forwardText, err := middleware(m, "")
		if err != nil {
			return
		}
		replyText := ""
		var atElement *message.AtElement
		if m.ReplyTo != nil {
			var err error
			replyText, err, atElement = getReplyText(m.ReplyTo)
			if err != nil {
				return
			}
		}
		messList := &message.SendingMessage{Elements: []message.IMessageElement{message.NewText(fmt.Sprintf("%s%s%s", replyText, forwardText, m.Text))}}
		if atElement != nil {
			messList.Append(atElement)
		}

		go bot.SendGroupMessage(QInfo.QQNumber, messList)
	})

	Bot.Handle(tb.OnDocument, func(m *tb.Message) {
		if m.Chat.Type != tb.ChatGroup && m.Chat.Type != tb.ChatSuperGroup {
			return
		}
		QInfo, ok := GroupFlag[m.Chat.ID]
		if !ok || (ok && !QInfo.Flag) {
			return
		}
		forwardText, err := middleware(m, "")
		if err != nil {
			return
		}
		file := m.Document
		fmt.Println(file.MIME)
		if strings.Contains(file.MIME, "image") {
			go func() {
				data, err := Bot.GetFile(&file.File)
				if err != nil {
					logger.Error(err)
				}

				var dataB []byte
				if strings.HasSuffix(file.File.FilePath, ".webp") {
					pic, err := webp.Decode(data, &decoder.Options{})
					if err != nil {
						logger.Error(err)
					}
					buf := new(bytes.Buffer)
					err = png.Encode(buf, pic)
					dataB = buf.Bytes()
				} else {
					dataB, err = ioutil.ReadAll(data)
					if err != nil {
						logger.Error(err)
					}
				}

				GroupImage, err := bot.UploadGroupImage(QInfo.QQNumber, bytes.NewReader(dataB))
				if err != nil {
					logger.Error(err)
				}
				go bot.SendGroupMessage(QInfo.QQNumber, &message.SendingMessage{Elements: []message.IMessageElement{message.NewText(fmt.Sprintf("%s发送图片 %s", forwardText, m.Caption)), GroupImage}})

			}()
		}

	})

	Bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		if m.Chat.Type != tb.ChatGroup && m.Chat.Type != tb.ChatSuperGroup {
			return
		}
		QInfo, ok := GroupFlag[m.Chat.ID]
		if !ok || (ok && !QInfo.Flag) {
			return
		}
		forwardText, err := middleware(m, "")
		if err != nil {
			return
		}
		replyText := ""
		var atElement *message.AtElement = nil

		if m.ReplyTo != nil {
			var err error
			replyText, err, atElement = getReplyText(m.ReplyTo)
			if err != nil {
				return
			}
		}

		go func() {
			data, err := Bot.GetFile(&m.Photo.File)
			if err != nil {
				logger.Error(err)
			}

			var dataB []byte
			if strings.HasSuffix(m.Photo.File.FilePath, ".webp") {
				pic, err := webp.Decode(data, &decoder.Options{})
				if err != nil {
					logger.Error(err)
				}
				buf := new(bytes.Buffer)
				err = png.Encode(buf, pic)
				dataB = buf.Bytes()
			} else {
				dataB, err = ioutil.ReadAll(data)
				if err != nil {
					logger.Error(err)
				}
			}

			GroupImage, err := bot.UploadGroupImage(QInfo.QQNumber, bytes.NewReader(dataB))
			if err != nil {
				logger.Error(err)
			}
			messList := &message.SendingMessage{Elements: []message.IMessageElement{message.NewText(fmt.Sprintf("%s%s发送图片 %s", replyText, forwardText, m.Caption)), GroupImage}}
			if atElement != nil {
				messList.Append(atElement)
			}

			bot.SendGroupMessage(QInfo.QQNumber, messList)
		}()
	})
	Bot.Handle(tb.OnSticker, func(m *tb.Message) {
		if m.Chat.Type != tb.ChatGroup && m.Chat.Type != tb.ChatSuperGroup {
			return
		}
		QInfo, ok := GroupFlag[m.Chat.ID]
		if !ok || (ok && !QInfo.Flag) {
			return
		}
		forwardText, err := middleware(m, "")
		if err != nil {
			return
		}
		replyText := ""
		var atElement *message.AtElement = nil

		if m.ReplyTo != nil {
			var err error
			replyText, err, atElement = getReplyText(m.ReplyTo)
			if err != nil {
				return
			}
		}
		data, err := Bot.GetFile(&m.Sticker.Thumbnail.File)
		emoji := m.Sticker.Emoji

		var dataB []byte
		if strings.HasSuffix(m.Sticker.Thumbnail.File.FilePath, ".webp") {
			pic, err := webp.Decode(data, &decoder.Options{})
			if err != nil {
				logger.Error(err)
			}
			buf := new(bytes.Buffer)
			err = png.Encode(buf, pic)
			dataB = buf.Bytes()
		} else {
			dataB, err = ioutil.ReadAll(data)
			if err != nil {
				logger.Error(err)
			}
		}
		GroupImage, err := bot.UploadGroupImage(QInfo.QQNumber, bytes.NewReader(dataB))
		if err != nil {
			logger.Error(err)
		}
		if err != nil {
			logger.Error(err)
		}
		dataB, err = ioutil.ReadAll(data)
		if err != nil {
			logger.Error(err)
		}
		messList := &message.SendingMessage{Elements: []message.IMessageElement{message.NewText(fmt.Sprintf("%s%s发送贴纸%s %s", replyText, forwardText, emoji, m.Caption)), GroupImage}}
		if atElement != nil {
			messList.Append(atElement)
		}
		go bot.SendGroupMessage(QInfo.QQNumber, messList)
		if m.Sticker.Animated && len(tgsAddr) > 0 {
			go bot.SendGroupMessage(QInfo.QQNumber, &message.SendingMessage{Elements: []message.IMessageElement{message.NewText(fmt.Sprintf("%s%s发送动态贴纸%s ,请在网页预览 %s/%s/preview", replyText, forwardText, emoji, tgsAddr, m.Sticker.FileID))}})
		}
	})
	go Bot.Start()

	bot.OnGroupMessage(func(qqClient *client.QQClient, groupMessage *message.GroupMessage) {

		tgID, ok := GroupMap[groupMessage.GroupCode]

		if ok {

			sender := groupMessage.Sender
			var content string
			var reply string
			var at string
			var atTGUser string
			var hasRich = false
			for _, ele := range groupMessage.Elements {

				switch o := ele.(type) {
				case *message.TextElement:

					content += o.Content

				case *message.ReplyElement:
					reply, atTGUser = parseReply(o, tgID.ID)
					println(groupMessage.ToString())
				case *message.AtElement:
					if reply != "" && at == "" {
						at = fmt.Sprintf("%s (%d)", o.Display, o.Target)
					} else {
						content += fmt.Sprintf("%s (%d)", o.Display, o.Target)
					}
				case *message.GroupImageElement:
					hasRich = true
					go func() {
						_, err := Bot.Send(tgID, &tb.Photo{Caption: fmt.Sprintf("%s (%d)☝", sender.DisplayName(), sender.Uin), File: tb.FromURL(o.Url)})
						if err != nil {
							logger.Error(err)
						}

					}()

				case *message.MusicShareElement:
					hasRich = true

					fmt.Println(o.Title)
					content += fmt.Sprintf("音乐分享: *%s* %s\n %s", o.Title, o.Brief, o.Url)

				case *message.LightAppElement:
					hasRich = true

					fmt.Println(o.Content)
					var app AppModel
					err := json.Unmarshal([]byte(strings.ReplaceAll(o.Content, "\"ver1\",", "")), &app)
					if err != nil {
						logger.Error(err)
					}
					switch app.View {
					case "music":
						musicInfo := app.Meta.Music
						fmt.Println(musicInfo.MusicURL)

						if strings.Contains(musicInfo.MusicURL, "music.163.com/song") {

							logger.Info("网易云分享")
							url := strings.ReplaceAll(strings.ReplaceAll(musicInfo.MusicURL, "http://", "https://"), "/song/me/", "/song/media/")

							go func() {
								resp, err := http.Get(url)
								if err != nil {
									logger.Error(err)
								}

								fmt.Println(url)
								_, err = Bot.Send(tgID, &tb.Audio{Title: musicInfo.Title, Performer: musicInfo.Desc, Caption: fmt.Sprintf("来自%s[分享自网易云音乐]", sender.DisplayName()), File: tb.FromReader(resp.Body), Thumbnail: &tb.Photo{File: tb.FromURL(musicInfo.Preview + "?param=300x300")}, MIME: "audio/mpeg"})
								if err != nil {
									logger.Error(err)
								}
							}()

						} else {
							logger.Info("其他分享")
							go func() {
								resp, err := http.Get(musicInfo.MusicURL)
								if err != nil {
									logger.Error(err)
								}
								_, err = Bot.Send(tgID, &tb.Audio{Title: musicInfo.Title, Performer: musicInfo.Desc, Caption: fmt.Sprintf("来自%s[分享自%s]", sender.DisplayName(), musicInfo.Tag), File: tb.FromReader(resp.Body), Thumbnail: &tb.Photo{File: tb.FromURL(musicInfo.Preview)}})
								if err != nil {
									logger.Error(err)
								}
							}()
						}
					case "news":
						newsInfo := app.Meta.News
						go func() {
							text := fmt.Sprintf("\\[来自%s\\] * %s *\n%s", MDReplace.Replace(newsInfo.Tag), MDReplace.Replace(newsInfo.Title), MDReplace.Replace(newsInfo.JumpURL))
							_, err := Bot.Send(tgID, fmt.Sprintf("* %s * \\(_%d_\\) : %s", MDReplace.Replace(sender.DisplayName()), sender.Uin, text), tb.ModeMarkdownV2)
							if err != nil {
								logger.Error(err)
							}
						}()
					default:
						if app.Meta.Detail1.Title != "" {
							info := app.Meta.Detail1
							text := fmt.Sprintf("[来自%s] * %s *\n%s", MDReplace.Replace(info.Title), MDReplace.Replace(info.Desc), MDReplace.Replace(info.Qqdocurl))
							_, err := Bot.Send(tgID, fmt.Sprintf("* %s * \\(_%d_\\) : %s", MDReplace.Replace(sender.DisplayName()), sender.Uin, text), tb.ModeMarkdownV2)
							if err != nil {
								logger.Error(err)
							}

						}

					}
					return

				}
			}
			println(sender, reply, at)
			//if sender.CardName != "" {
			//	senderStr=sender.CardName
			//}else {
			//	senderStr=sender.DisplayName()
			//}
			if len(content) > 0 {
				go func() {
					if len(reply) > 0 {
						text := fmt.Sprintf("\\> *_ %s _* : *_ %s _*\n __\\-\\-\\-\\-\\-\\-\\-\\-__ \n * %s * \\(_%d_\\) : %s %s", MDReplace.Replace(at), MDReplace.Replace(reply), MDReplace.Replace(sender.DisplayName()), sender.Uin, MDReplace.Replace(content), atTGUser)

						_, err := Bot.Send(tgID, text, tb.ModeMarkdownV2)
						if err != nil {
							logger.Error(err)
						}
					} else {
						_, err := Bot.Send(tgID, fmt.Sprintf("* %s * \\(_%d_\\) : %s %s", MDReplace.Replace(sender.DisplayName()), sender.Uin, MDReplace.Replace(groupMessage.ToString()), atTGUser), tb.ModeMarkdownV2)
						if err != nil {
							logger.Error(err)
						}
					}

				}()

			} else if hasRich == false {
				_, err := Bot.Send(tgID, fmt.Sprintf("* %s * \\(_%d_\\) : %s %s", MDReplace.Replace(sender.DisplayName()), sender.Uin, MDReplace.Replace(groupMessage.ToString()), atTGUser), tb.ModeMarkdownV2)
				if err != nil {
					logger.Error(err)
				}
			}
		}

	})
}

func parseReply(m *message.ReplyElement, tgID int64) (string, string) {
	content := ""
	atTGUser := ""
	if m.Sender == config.GlobalConfig.GetInt64("bot.account") {
		var element *message.TextElement
		for _, ele := range m.Elements {
			if textElement, ok := ele.(*message.TextElement); ok {
				element = textElement
				break
			}
		}
		if element != nil {

			res := strings.Split(element.Content, "-------  ")
			var text = res[len(res)-1]

			names := strings.Split(text, " :")
			name := names[0]

			member, err := GetAMember(tgID, strings.TrimSpace(name))
			if err == nil {
				atTGUser = fmt.Sprintf("[%s](tg://user?id=%d)", MDReplace.Replace(name), member.UserID)
			}
		}

	}

	for _, ele := range m.Elements {
		switch o := ele.(type) {
		case *message.TextElement:
			content += o.Content

		}

	}

	return content, atTGUser
}

func (m *tgForward) Start(bot *bot.Bot) {
	//panic("implement me")
}

func (m *tgForward) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	//panic("implement me")
}

var instance *tgForward
var Bot *tb.Bot

func init() {
	instance = &tgForward{}
	bot.RegisterModule(instance)
}
func (m *tgForward) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "external.TGForward",
		Instance: instance,
	}
}
