package sub

import (
	"fmt"
	"formbot/cmd/form/subform/function"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleAddCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// 回答が正しく得られなかった場合，終了
	if len(options) != 2 {
		log.Printf("invalid options: %#v", options)
		return
	}

	channelID := i.ChannelID
	channelName := "ChannelName"
	hour := options[0].IntValue()
	day := options[1].StringValue()

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

	// 時間が正しく入力されなかった場合，終了
	if hour < 0 || 23 < hour {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "時間は0~23の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}


	// ファイルへの書き込み
	err = subfunc.WriteToDataFile(channelID, channelName, fmt.Sprintf("%d", hour), day)
	if err != nil {
		log.Printf("failed to write data to file: %v", err)
		return
	}

	// アラームセット完了報告
	dayJ, err := subfunc.WeekEtoJ(day)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
		return
	}
	name := i.Member.Nick
	if name == "" {
		name = i.Member.User.Username
	}
	mes := discordgo.MessageEmbed{
		Color:       0xF1C40F,
		Footer:      &discordgo.MessageEmbedFooter{Text: "通知が設定されました."},
		Description: fmt.Sprintf("毎週%s曜日の%d時にリマインドを通知します.", dayJ, hour),
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

	// 通知追加中の表示を変更する
	finishFollowUpStr := "Addコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
	
}
