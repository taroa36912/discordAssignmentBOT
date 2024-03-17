package add

import (
	"fmt"
	"formbot/send/send_e"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleAddCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	channelID := i.ChannelID
	channelName := "ChannelName" // 本当のチャンネル名をここに設定
	hour := options[0].IntValue()
	day := options[1].StringValue()

	// ファイルへの書き込み
	err := form.WriteToDataFile(channelID, channelName, fmt.Sprintf("%d", hour), day)
	if err != nil {
		log.Printf("failed to write data to file: %v", err)
		return
	}

	// アラームセット完了報告
	dayJ, err := form.WeekEtoJ(day)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
		return
	}

	result := fmt.Sprintf("毎週%s曜日の%d時リマインドを通知します.", dayJ, hour)
	send.SendMessage(s, i.ChannelID, result)
}