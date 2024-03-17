package send

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"os"
)

const (
	EnvClientId = "client_id"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv(EnvClientId)
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
	// 個々の部分を書き換える
	if u.ID != clientId {
		subfunc.SendReply(s, m.ChannelID,"ok", m.Reference())
	}
}
