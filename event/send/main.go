package send

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
	// メッセージが!cで始まる場合の処理
	if strings.HasPrefix(m.Content, "!c") {
		// ここに!cで始まるメッセージの処理を書く
		subfunc.SendMessage(s, m.ChannelID, "!cが使われました.")
	}
}
