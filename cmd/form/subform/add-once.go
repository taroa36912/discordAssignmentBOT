package sub

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)



func HandleAddOnceCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// 回答が正しく得られなかった場合，終了
	if len(options) != 3 {
		log.Printf("invalid options: %#v", options)
		return
	}

	channelID := i.ChannelID
	dateString := options[0].StringValue()
	channelName := options[1].StringValue()
	mention := options[2].StringValue()
	// time.Parse を使って時間を解析し、変数に格納
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		fmt.Println("Failed to parse time:", err)
		return
	}

	// 年、月、日、時を取得して変数に格納
	year := parsedTime.Year()
	month := parsedTime.Month()
	day := parsedTime.Day()
	hour := parsedTime.Hour()

	// 処理を行っている間表示されるメッセージ
	followUp := discordgo.WebhookParams{
		Content: "通知追加中...",
		Flags:   discordgo.MessageFlagsEphemeral,
	}
	followUpMsg, err := s.FollowupMessageCreate(i.Interaction, true, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}

	// ファイルへの書き込み
	if(mention == "me"){
		// 自分のみならば，DMを送信する
		channel, err := s.UserChannelCreate(i.Member.User.ID) // DMの生成
		if err != nil {
			fmt.Println("Error creating DM channel: ", err)
			return
		}
		err = subfunc.WritetoFile("form.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %d, %s", channel.ID, channelName, year, month, day, hour, mention))
		if err != nil {
			log.Printf("failed to write data to file: %v", err)
			return
		}
		// メッセージを送信
		_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("```あなただけに見えるリマインドを設定しました.\n当日 : %s, %d年%d月%d日%d時にお知らせします.```", channelName, year, month, day, hour))
		if err != nil {
			fmt.Println("Error sending message: ", err)
			return
		}
	}else{
		err = subfunc.WritetoFile("form.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %d, %s", channelID, channelName, year, month, day, hour, mention))
		if err != nil {
			log.Printf("failed to write data to file: %v", err)
			return
		}
		name := i.Member.Nick
		if name == "" {
			name = i.Member.User.Username
		}
		mes := discordgo.MessageEmbed{
			Color:       0x800020,
			Footer:      &discordgo.MessageEmbedFooter{Text: "通知が設定されました."},
			Description: fmt.Sprintf("%d年%d月%d日%d時にリマインドを通知します.", year, month, day, hour),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    name,
				URL:     fmt.Sprintf("https://discordapp.com/users/%s", i.Member.User.ID),
				IconURL: i.Member.User.AvatarURL(""),
			},
		}
		_, err = s.ChannelMessageSendEmbed(i.ChannelID, &mes)
		if err != nil {
			log.Printf("failed to send message, err: %v", err)
		}
	}

	// 通知追加中の表示を変更する
	finishFollowUpStr := "Add-onceコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
}