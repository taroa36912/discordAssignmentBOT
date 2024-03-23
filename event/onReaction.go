package event

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"formbot/function"
	"log"
)

func MessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 自分のみのチャンネルに，通知一覧を送信する
	zemiData, err := subfunc.ReadZemiFile()
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, data := range zemiData {

	}
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
