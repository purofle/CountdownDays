package bot

import "github.com/mymmrac/telego"

func PrivateChatOnly() func(update telego.Update) bool {
	return func(update telego.Update) bool {
		return update.Message.Chat.Type == telego.ChatTypePrivate && !update.Message.From.IsBot
	}
}
