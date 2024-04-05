package event

import (
    "github.com/bwmarrin/discordgo"
)

// MaxMessageLength は、discordメッセージの最大文字数です。
const MaxMessageLength = 2000

func OverWriteMessage(s *discordgo.Session, m *discordgo.MessageCreate){
    // メッセージを分割して送信する
    for _, part := range splitMessage(m.Content, MaxMessageLength) {
    	s.ChannelMessageSend(m.ChannelID, part)
    }
}

// splitMessage は、指定された文字数でメッセージを分割します。
func splitMessage(message string, maxLen int) []string {
    var parts []string
    for len(message) > maxLen {
        part := message[:maxLen]
        parts = append(parts, part)
        message = message[maxLen:]
    }
    parts = append(parts, message)
    return parts
}
