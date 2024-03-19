package checkform

import (
	"fmt"
	"formbot/function"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CheckReminder(s *discordgo.Session, e *discordgo.Ready){
	for {
		// 現在の時刻と曜日を取得
		current := time.Now()
		remindData, err := subfunc.ReadDataFile()
		if err != nil {
			log.Printf("failed to get data.txt: %v", err)
			return
		}
		// ここに定期的に行いたい処理を記述する
		fmt.Println("Periodic task executed at: ", time.Now())
		for _, data := range remindData {
			channelID, sentence := subfunc.CheckEachRow(s, e, data, current)
			if channelID != "" && sentence != ""{
				subfunc.SendMessage(s, channelID, sentence)
			}
		}

		// 1時間待機する
		time.Sleep(time.Hour)
	}

}