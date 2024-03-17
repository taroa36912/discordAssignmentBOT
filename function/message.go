package subfunc

import (
	"github.com/bwmarrin/discordgo"
	"log"
)


func SendMessage(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message : ", err)
	}
}