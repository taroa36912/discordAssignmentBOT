package send

import (
	"fmt"
	"formbot/send/send_e"
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
	if u.ID != clientId {
		send.SendMessage(s, m.ChannelID, u.Mention()+"ok1")
		send.SendReply(s, m.ChannelID,"ok2", m.Reference())
	}
}
