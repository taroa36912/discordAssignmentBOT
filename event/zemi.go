package event

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"time"
)

const (
	EnvZemiChannel = "zemi_channel_id"
)


func CreateZemiMessage(s *discordgo.Session, e *discordgo.Ready) {

}



func CheckZemiReaction(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv(EnvZemiChannel)
	attendanceEmoji := ":syusseki:"
	// absentEmoji := ":kesseki:"
	// nakairiEmoji := ":totyusanka:"
	// nakanukeEmoji := ":totyunuke:"
	// 定期的な絵文字チェックの開始
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		// zemi出席のメッセージIDそれぞれに対し処理
		zemiData, err := subfunc.ReadFile("zemiMessage.txt")
		if err != nil {
			log.Printf("failed to get data.txt: %v", err)
			return
		}
		for _, messageID := range zemiData {
			// ここに処理を書く
		}
	}
}
