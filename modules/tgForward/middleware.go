package tgForward

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
)

func fwUserName(m *tb.Message, forwardText string) string {
	return fmt.Sprintf("%s %s %s :", forwardText, m.Sender.FirstName, m.Sender.LastName)
}
func fwTGSource(m *tb.Message, forwardText string) string {
	if m.IsForwarded() {
		sourceText := ""
		if m.OriginalChat != nil {
			sourceText = m.OriginalChat.Title
		}
		if m.OriginalSenderName != "" {
			sourceText += ":" + m.OriginalSenderName
		}

		return fmt.Sprintf("%s [From %s]", forwardText, sourceText)

	}
	return forwardText
}
