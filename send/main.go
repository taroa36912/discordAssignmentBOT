package send

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"formbot/send/send_e"
)

const (
	EnvClientId     = "client_id"
)


func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv(EnvClientId)
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
	if u.ID != clientId {
		send.SendMessage(s, m.ChannelID, u.Mention() + m.Content)
		send.SendReply(s, m.ChannelID, m.Content, m.Reference())
	}
}