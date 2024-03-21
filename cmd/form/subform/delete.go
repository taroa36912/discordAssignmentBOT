package sub

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func HandleDeleteCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {

	// 回答が正しく得られなかった場合，終了
	if len(options) != 1 {
		log.Printf("invalid options: %#v", options)
		return
	}

	index := options[0].IntValue()
	
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
	// 通知番号の表示
	count := 1
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	for _, data := range remindData {
		sentence := subfunc.ViewEachRow(channel.ID, data)
		if sentence != "" {
			if count == int(index){
				mention := subfunc.MentionType(data)
				if mention == "everyone"{subfunc.SendMessage(s, strings.Split(data, ", ")[0], fmt.Sprintf("@everyone```以下の通知を削除します.\n%s```", sentence))}else if mention == "me"{subfunc.SendMessage(s, strings.Split(data, ", ")[0], fmt.Sprintf("```以下の通知を削除します.\n%s```", sentence))}
				subfunc.ReadAndDeleteDataFile(data)
			}
			count++
		}
	}

	// 表示を変更する
	finishFollowUpStr := "deleteコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
}
