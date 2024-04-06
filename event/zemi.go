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


// 毎日20時にゼミの出席メッセージを作る関数
func CreateZemiMessage(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv("zemi_channel_id")
	zemiRoleID := os.Getenv("zemi_role_id")
	zemiWeek := "Tuesday"
	zemiHour := 10
	zemiMinute := 0
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 一週間後の時刻を取得
	current := time.Now().In(location)
	oneWeekLater := current.AddDate(0, 0, 7)
	currentweekday := current.Weekday().String()
	if currentweekday == zemiWeek {
		year, month, day := oneWeekLater.Date()
		dayJ, err := subfunc.WeekEtoJ(zemiWeek)
		if err != nil {
			log.Printf("failed to convert day to Japanese: %v", err)
			return
		}
		sentence := fmt.Sprintf("<@&%s>```%d年%d月%d日%s曜日%d時%d分~\n自主ゼミの出欠を取ります.\nリアクションをしてください.```", zemiRoleID, year, month, day, dayJ, zemiHour, zemiMinute)
		msg, err := s.ChannelMessageSend(zemiChannelID, sentence)
		if err != nil {
			log.Println("Error sending message : ", err)
			return
		}
		subfunc.WritetoFile("zemiMessage.txt", fmt.Sprintf("%s, %d, %d, %d, %s, %d, %d", msg.ID, year, month, day, zemiWeek, zemiHour, zemiMinute))
	}
}


// 毎朝8時に呼び，その日に開催するゼミのリアクションをチェックする関数
func CheckZemiReaction(s *discordgo.Session, e *discordgo.Ready) {
	zemiServerID := os.Getenv("guild_id")
	zemiChannelID := os.Getenv("zemi_channel_id")
	zemiRoleID := os.Getenv("zemi_role_id")
	zemiAttendID := os.Getenv("emoji_attend_id")
	zemiAbsentID := os.Getenv("emoji_absent_id")
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

	// 現在の時刻と曜日を取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 現在時刻を取得
	current := time.Now().In(location)

	for _, message := range zemiMessage {
		parts := strings.Split(message, ", ")
		if len(parts) == 7 {
			messageID := parts[0]
			year := parts[1]
			month := parts[2]
			day := parts[3]
			weekday := parts[4]
			hour := parts[5]
			minute := parts[6]
			dayJ, err := subfunc.WeekEtoJ(weekday)
			if err != nil {
				log.Printf("failed to convert day to Japanese: %v", err)
				return
			}
			currentYear := current.Year()
			currentMonth := int(current.Month())
			currentday := current.Day()

			if(fmt.Sprintf("%d", currentYear) == year && fmt.Sprintf("%d", currentMonth) == month && fmt.Sprintf("%d", currentday) == day){
				// メッセージを取得
				_, err = s.ChannelMessage(zemiChannelID, messageID)
				if err != nil {
					s.ChannelMessageSend(zemiChannelID, "ゼミがキャンセルされました.")
					subfunc.DeleteFile("zemiMessage.txt", message)
					return
				}

				// 出席リアクションをしたユーザーを取得
				usersAttend, err := s.MessageReactions(zemiChannelID, messageID, zemiAttendID, 100, "", "")
				if err != nil {
					fmt.Println("Error getting reactions: ", err)
					return
				}
				usersMentionAttend := usersMention(usersAttend)

				// 欠席リアクションをしたユーザーを取得
				usersAbsent, err := s.MessageReactions(zemiChannelID, messageID, zemiAbsentID, 100, "", "")
				if err != nil {
					fmt.Println("Error getting reactions: ", err)
					return
				}
				usersMentionAbsence := usersMention(usersAbsent)


				// 返信先のメッセージの参照情報
				reference := &discordgo.MessageReference{
					MessageID: messageID,
					ChannelID: zemiChannelID,
					GuildID: zemiServerID,
				}
				sentence := fmt.Sprintf("<@&%s>```%s年%s月%s日%s曜日%s時%s分\n自主ゼミ当日です.\n```参加者 : %s\n欠席者 : %s", zemiRoleID, year, month, day, dayJ, hour, minute, usersMentionAttend, usersMentionAbsence)
				// SendReply関数を呼び出してメッセージを送信
				_, err = s.ChannelMessageSendReply(zemiChannelID, sentence, reference)
				if err != nil {
					log.Println("Error sending message: ", err)
				}
			}
		}
	}
}


// ユーザーリストから，ユーザーIDを結合する関数
func usersMention(users []*discordgo.User)(string){
	userslist := ""
	for _, user := range users{
		userslist += "<@" + user.ID + ">, "
	}
	return userslist
}


// 毎時0分に呼び出し，その日に開始されるゼミの時間に一致するものを通知する関数
func ZemiTimeNotification(s *discordgo.Session, e *discordgo.Ready) {
	zemiChannelID := os.Getenv("zemi_channel_id")
	zemiRoleID := os.Getenv("zemi_role_id")
	// zemi出席のメッセージIDそれぞれに対し処理
	zemiMessage, err := subfunc.ReadFile("zemiMessage.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}

	// 現在の時刻と曜日を取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 現在時刻を取得
	current := time.Now().In(location)

	for _, message := range zemiMessage {
		parts := strings.Split(message, ", ")
		if len(parts) == 7 {
			messageID := parts[0]
			year := parts[1]
			month := parts[2]
			day := parts[3]
			weekday := parts[4]
			hour := parts[5]
			minute := parts[6]
			currentYear := current.Year()
			currentMonth := int(current.Month())
			currentday := current.Day()
			currentHour := current.Hour()
			if fmt.Sprintf("%d", currentYear) == year && fmt.Sprintf("%d", currentMonth) == month && fmt.Sprintf("%d", currentday) == day && fmt.Sprintf("%d", currentHour) == hour{
				dayJ, err := subfunc.WeekEtoJ(weekday)
				if err != nil {
					log.Printf("failed to convert day to Japanese: %v", err)
					return
				}
				// 続きの処理
				// メッセージを取得
				_, err = s.ChannelMessage(zemiChannelID, messageID)
				if err != nil {
					s.ChannelMessageSend(zemiChannelID, "ゼミがキャンセルされました.")
					subfunc.DeleteFile("zemiMessage.txt", message)
					return
				}

				// 返信先のメッセージの参照情報
				reference := &discordgo.MessageReference{
					MessageID: messageID,
				}
				sentence := fmt.Sprintf("<@%s>```%s年%s月%s日%s曜日%s時%s分\n自主ゼミが開始されます.```", zemiRoleID, year, month, day, dayJ, hour, minute)
				// SendReply関数を呼び出してメッセージを送信
				_, err = s.ChannelMessageSendReply(zemiChannelID, sentence, reference)
				if err != nil {
					log.Println("Error sending message: ", err)
				}
				subfunc.DeleteFile("zemiMessage.txt", message)
			}
		}
	}
}
