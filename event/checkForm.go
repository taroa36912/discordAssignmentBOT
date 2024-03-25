package event

import (
	"formbot/function"
	"log"
	"time"
    "fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CheckReminder(s *discordgo.Session, e *discordgo.Ready) {
    // 現在の時刻と曜日を取得
    // 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 現在時刻を取得
	current := time.Now().In(location)
    remindData, err := subfunc.ReadFile("form.txt")
    if err != nil {
        log.Printf("failed to get data.txt: %v", err)
        return
    }
    for _, data := range remindData {
        channelID, sentence := checkEachRow(s, e, data, current)
        if channelID != "" && sentence != "" {
            // メッセージを送信
	        _, err = s.ChannelMessageSend(channelID, sentence)
		    if err != nil {
		        fmt.Println("Error sending message: ", err)
		        return
	        }
        }
    }
}



func checkEachRow(s *discordgo.Session, e *discordgo.Ready, data string, current time.Time)(string, string){
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ5はweekly
	if len(parts) == 5 {
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		hour := parts[2]
		weekday := parts[3]
		mention := parts[4]

		currentHour := current.Hour()
		currentWeekday := current.Weekday().String()

		// 記録された時刻と曜日と現在の時刻と曜日が一致する場合にSendMessageを実行
		if hour == fmt.Sprintf("%02d", currentHour) && weekday == currentWeekday {
			if(mention == "everyone"){
				sentence := fmt.Sprintf("@everyone\n```課題：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			}else if(mention == "me"){
				sentence := fmt.Sprintf("```課題：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			}
		}
	}else if len(parts) == 7 { // 長さ7はonce
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		year := parts[2]
		month := parts[3]
		day := parts[4]
		hour := parts[5]
		mention := parts[6]

		currentHour := current.Hour()
		currentYear := current.Year()
		currentMonth := int(current.Month())
		currentDay := current.Day()

		// 指定された年月日時が現在時刻と一致する場合にのみ処理を実行
		if year == fmt.Sprintf("%d", currentYear) &&
			month == fmt.Sprintf("%d", currentMonth) &&
			day == fmt.Sprintf("%d", currentDay) &&
			hour == fmt.Sprintf("%d", currentHour) {
			subfunc.DeleteFile("form.txt", data)
			if mention == "everyone" {
				sentence := fmt.Sprintf("@everyone\n```当日：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			} else if mention == "me" {
				sentence := fmt.Sprintf("```当日：%sの締め切りを通知します.```", title)
				return channelID, sentence
			}
		}
	}
	return "", ""
}