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
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 現在時刻を取得
	current := time.Now().In(location)

	//毎晩20時に実行
	zemiData, err := subfunc.ReadFile("zemiData.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	oneWeekLater := current.AddDate(0, 0, 7)
	for _, zemi := range zemiData {
		// ここに処理を書く
		parts := strings.Split(zemi, ", ")
		if len(parts) == 2 {
			mentionID := parts[0]
			weekday := parts[1]
			teacher, err := subfunc.GetTeacherName(mentionID)
			if err != nil {
				log.Printf("failed to get teacher name: %v", err)
				return
			}
			dayJ, err := subfunc.WeekEtoJ(weekday)
			if err != nil {
				log.Printf("failed to convert day to Japanese: %v", err)
				return
			}
			year, month, day := oneWeekLater.Date()
			currentweekday := current.Weekday().String()
			if weekday == currentweekday {
				sentence := fmt.Sprintf("<@%s>```%d年%d月%d日%s曜日\n%sゼミの出欠を取ります.\nリアクションをしてください.```", mentionID, year, month, day, dayJ, teacher)
				msg, err := s.ChannelMessageSend(zemiChannelID, sentence)
				if err != nil {
					log.Println("Error sending message : ", err)
				}
				subfunc.WriteFile("zemiMessage.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %s", msg.ID, mentionID, year, month, day, weekday))
			}
		}
	}
}

// 8時に呼び，リアクションをチェックする関数
func CheckZemiReaction(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv(EnvZemiChannel)
	// attendanceEmoji := ":syusseki:"
	// absentEmoji := ":kesseki:"
	// nakairiEmoji := ":totyusanka:"
	// nakanukeEmoji := ":totyunuke:"
	// 定期的な絵文字チェックの開始
	// zemi出席のメッセージIDそれぞれに対し処理
	zemiMessage, err := subfunc.ReadFile("zemiMessage.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, message := range zemiMessage {
		parts := strings.Split(message, ", ")
		if len(parts) == 6 {
			messageID := parts[0]
			mentionID := parts[1]
			year := parts[2]
			month := parts[3]
			day := parts[4]
			weekday := parts[5]
			teacher, err := subfunc.GetTeacherName(mentionID)
			if err != nil {
				log.Printf("failed to get teacher name: %v", err)
				return
			}
			dayJ, err := subfunc.WeekEtoJ(weekday)
			if err != nil {
				log.Printf("failed to convert day to Japanese: %v", err)
				return
			}
			// 続きの処理
			// メッセージを取得
			_, err = s.ChannelMessage(zemiChannelID, messageID)
			if err != nil {
				subfunc.DeleteFile("zemiMessage.txt", message)
				return
			}

			// 返信先のメッセージの参照情報
			reference := &discordgo.MessageReference{
				MessageID: messageID,
			}
			sentence := fmt.Sprintf("<@%s>```%s年%s月%s日%s曜日\n%sゼミ当日です.\nリアクションをしてください.```", mentionID, year, month, day, dayJ, teacher)
			// SendReply関数を呼び出してメッセージを送信
			_, err = s.ChannelMessageSendReply(zemiChannelID, sentence, reference)
			if err != nil {
				log.Println("Error sending message: ", err)
			}
		}
	}
}

// ゼミの時間直前に呼び出し，ゼミの時間に一致するものを通知する関数
func ZemiTimeNotification(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv(EnvZemiChannel)
	// zemi出席のメッセージIDそれぞれに対し処理
	zemiMessage, err := subfunc.ReadFile("zemiMessage.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, message := range zemiMessage {
		parts := strings.Split(message, ", ")
		if len(parts) == 6 {
			messageID := parts[0]
			mentionID := parts[1]
			year := parts[2]
			month := parts[3]
			day := parts[4]
			weekday := parts[5]
			// 現在の時刻と曜日を取得
			// 日本標準時の場所情報を取得
			location, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				fmt.Println("Failed to load location:", err)
				return
			}
			// 現在時刻を取得
			current := time.Now().In(location)
			currentYear := current.Year()
			currentMonth := int(current.Month())
			currentday := current.Day()
			if year == currentYear && month == currentMonth && day == currentday {
				teacher, err := subfunc.GetTeacherName(mentionID)
				if err != nil {
					log.Printf("failed to get teacher name: %v", err)
					return
				}
				dayJ, err := subfunc.WeekEtoJ(weekday)
				if err != nil {
					log.Printf("failed to convert day to Japanese: %v", err)
					return
				}
				// 続きの処理
				// メッセージを取得
				_, err = s.ChannelMessage(zemiChannelID, messageID)
				if err != nil {
					s.ChannelMessageSend(zemiChannelID, fmt.Sprintf("%sゼミがキャンセルされました.", teacher))
					subfunc.DeleteFile("zemiMessage.txt", message)
					return
				}

				// 返信先のメッセージの参照情報
				reference := &discordgo.MessageReference{
					MessageID: messageID,
				}
				sentence := fmt.Sprintf("<@%s>```%s年%s月%s日%s曜日\n%sゼミが開始されます.```", mentionID, year, month, day, dayJ, teacher)
				// SendReply関数を呼び出してメッセージを送信
				_, err = s.ChannelMessageSendReply(zemiChannelID, sentence, reference)
				if err != nil {
					log.Println("Error sending message: ", err)
				}
			}
		}
	}
}
