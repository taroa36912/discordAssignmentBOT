package sub

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
)



func HandleAddOnceCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// 回答が正しく得られなかった場合，終了
	if len(options) != 6 {
		log.Printf("invalid options: %#v", options)
		return
	}

	year := options[0].IntValue()
	month := options[1].IntValue()
	day := options[2].IntValue()
	hour := options[3].IntValue()
	channelName := options[4].StringValue()
	mention := options[5].StringValue()
	// ロールIDからロール名を取得
	mentionName, err := subfunc.GetRoleName(s, i.GuildID, mention) // ここに対象のロールIDを入れます
	if err != nil {
		log.Println("ロール名の取得に失敗しました:", err)
		mentionName = "undefined"
	}
	channelID := i.ChannelID

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
		err = subfunc.WritetoFile("form.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %d, %s, %s", channel.ID, channelName, year, month, day, hour, mention, mentionName))
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
		err = subfunc.WritetoFile("form.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %d, %s, %s", channelID, channelName, year, month, day, hour, mention, mentionName))
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
			Description: fmt.Sprintf("メンション範囲 : %s\nメッセージ : %s\n%d年%d月%d日%d時にリマインドを通知します.", mentionName, channelName, year, month, day, hour),
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