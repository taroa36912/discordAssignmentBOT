package send

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := m.Author
	channelID := m.ChannelID
	message := m.Message
	// 取得したメッセージを表示
	fmt.Printf("%s %s(%s) > %s\n", channelID, u.Username, u.ID, message.Content)
}
