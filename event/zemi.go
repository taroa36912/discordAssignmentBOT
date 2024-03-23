package event

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
	"time"
)

const (
	EnvZemiChannel = "zemi_channel_id"
)

func CreateZemiMessage(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv(EnvZemiChannel)
	zemiData, err := subfunc.ReadFile("zemiData.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, zemi := range zemiData {
		// ここに処理を書く
		parts := strings.Split(zemi, ", ")
		if len(parts) == 4 {
			mentionID := parts[0]
			teacher := parts[1]
			weekday := parts[2]
			hour := parts[3]
			dayJ, err := subfunc.WeekEtoJ(weekday)
			if err != nil {
				log.Printf("failed to convert day to Japanese: %v", err)
				return
			}
			sentence := fmt.Sprintf("<@%s>```%s曜日 : %s時\n%s研のゼミの出席を取ります.\nリアクションをしてください.```", mentionID, dayJ, hour, teacher)
			msg, err := s.ChannelMessageSend(zemiChannelID, sentence)
			if err != nil {
				log.Println("Error sending message : ", err)
			}

			subfunc.WriteToDataFile("zemiMessage.txt", fmt.Sprintf("%s, %s, %s", msg.ID, mentionID, date))
		}
	}

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
