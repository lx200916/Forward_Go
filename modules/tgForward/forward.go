package tgForward

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"sync"
	"time"
)

type tgForward struct {
	
}
var logger = utils.GetModuleLogger("external.TGForward")
var GroupMap = make(map[int64]int64)

type groupConfig struct {
	QQ int64 `yaml:"QQ"`
	TG int64 `yaml:"TG"`
}
func (m *tgForward) Init() {
	var err error
	Bot, err = tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: config.GlobalConfig.GetString("Telegram.APIAddr"),
		Token:  config.GlobalConfig.GetString("Telegram.token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	var grouplists []groupConfig
	err=config.GlobalConfig.UnmarshalKey("Groups",&grouplists)
	if err != nil {
		logger.Error("Group Lists Fail to Parse")
	}
	for _,groupInfo:=range grouplists{
		GroupMap[groupInfo.QQ]=groupInfo.TG
	}
	
}

func (m *tgForward) PostInit() {
logger.Info("Init Finish")}

func (m *tgForward) Serve(bot *bot.Bot) {
	//panic("implement me")
	bot.OnGroupMessage(func(qqClient *client.QQClient, groupMessage *message.GroupMessage) {
		tgID,ok:=GroupMap[groupMessage.GroupCode]
		if ok {
			
		}

	})
}
func forwardText(b *bot.Bot)  {

}

func (m *tgForward) Start(bot *bot.Bot) {
	panic("implement me")
}

func (m *tgForward) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	panic("implement me")
}

var instance *tgForward
var Bot *tb.Bot
func init() {
	instance = &tgForward{}
	bot.RegisterModule(instance)
}
func (m *tgForward) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "internal.logging",
		Instance: instance,
	}
}