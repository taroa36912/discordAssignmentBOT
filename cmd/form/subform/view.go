package sub

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleViewCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// 処理を行っている間表示されるメッセージ
	followUp := discordgo.WebhookParams{
		Content: "通知一覧準備中...",
		Flags:   discordgo.MessageFlagsEphemeral,
	}
	followUpMsg, err := s.FollowupMessageCreate(i.Interaction, true, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}

	// ファイルの読み込み
	// 自分のみならば，DMを送信する
	channel, err := s.UserChannelCreate(i.Member.User.ID) // DMの生成
	if err != nil {
		fmt.Println("Error creating DM channel: ", err)
		return
	}

	// 自分のみのチャンネルに，通知一覧を送信する
	remindData, err := subfunc.ReadDataFile()
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, data := range remindData {
		sentence := subfunc.ViewEachRow(channel.ID, data)
		if sentence != ""{
			subfunc.SendMessage(s, channel.ID, sentence)
		}
	}
		
	// 表示を変更する
	finishFollowUpStr := "viewコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
}