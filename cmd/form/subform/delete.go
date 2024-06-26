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
		Content: "通知削除中...",
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
	remindData, err := subfunc.ReadFile("form.txt")
	// 通知番号の表示
	count := 1
	if err != nil {
		log.Printf("failed to get form.txt: %v", err)
		return
	}
	for _, data := range remindData {
		sentence := subfunc.ViewEachRow(channel.ID, data)
		if sentence != "" {
			if count == int(index){
				mention := subfunc.MentionType(data)
				if mention == "me"{
					// メッセージを送信
					_, err = s.ChannelMessageSend(strings.Split(data, ", ")[0], fmt.Sprintf("```以下の通知を削除します.\n%s```", sentence))
					if err != nil {
						fmt.Println("Error sending message: ", err)
						return
					}
				}else{
					// メッセージを送信
					_, err = s.ChannelMessageSend(strings.Split(data, ", ")[0], fmt.Sprintf("<@&%s>```以下の通知を削除します.\n%s```", mention, sentence))
					if err != nil {
						fmt.Println("Error sending message: ", err)
						return
					}
				}
				subfunc.DeleteFile("form.txt", data)
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
