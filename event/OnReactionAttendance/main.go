package reaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 特定のチャンネルとメッセージ、特定のリアクションでのみ処理を行う
	if r.ChannelID != "YOUR_CHANNEL_ID" || r.MessageID != "TARGET_MESSAGE_ID" || r.Emoji.Name != "👍" {
		return
	}

	// ユーザー情報を取得
	user, err := s.User(r.UserID)
	if err != nil {
		fmt.Println("Error getting user information:", err)
		return
	}

	// 出力
	fmt.Printf("Reaction added by %s#%s (ID: %s) to message ID: %s\n", user.Username, user.Discriminator, user.ID, r.MessageID)
}
