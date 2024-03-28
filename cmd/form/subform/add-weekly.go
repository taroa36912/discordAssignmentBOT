package sub

import (
	"fmt"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func HandleAddWeeklyCommand(
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
	// コマンドが投下されたチャンネルのチャンネル名を取得
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
    	log.Printf("failed to get channel information: %v", err)
    	return
	}

	channelName := channel.Name
	hour := options[0].IntValue()
	day := options[1].StringValue()
	mention := options[2].StringValue()

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

	// everyoneメンションが2つ以上あるとき, キャンセルして終了する
	remindData, err := subfunc.ReadFile("form.txt")
	if err != nil {
		log.Printf("failed to get data.txt: %v", err)
		return
	}
	count := 0
	for _, data := range remindData {
		flag := checkWeeklyEachRow(channelID, channelName, data)
		if flag{count++}
	}
	if count > 1 {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "2つ以上のweekly-everyone通知は設定できません."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
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

	// アラームセット完了報告
	dayJ, err := subfunc.WeekEtoJ(day)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
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
		err = subfunc.WritetoFile("form.txt", fmt.Sprintf("%s, %s, %d, %s, %s", channel.ID, channelName, hour, day, mention))
		if err != nil {
			log.Printf("failed to write data to file: %v", err)
			return
		}
		// メッセージを送信
		_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("```あなただけに見えるリマインドを設定しました.\n当日 : %s, 毎週%s曜日の%d時にお知らせします.```", channelName, dayJ, hour))
		if err != nil {
			fmt.Println("Error sending message: ", err)
			return
		}
	}else{
		err = subfunc.WritetoFile("data.txt", fmt.Sprintf("%s, %s, %d, %s, %s", channelID, channelName, hour, day, mention))
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
	}

	// 通知追加中の表示を変更する
	finishFollowUpStr := "Add-weeklyコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
	
}



func checkWeeklyEachRow(defaultChannelID string, defaltChannelName string, data string)(bool){
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ5はweekly
	if len(parts) == 5 {
		channelID := parts[0]
		title := parts[1]
		mention := parts[4]
		// 記録された時刻と曜日と現在の時刻と曜日が一致する場合にSendMessageを実行
		if(defaultChannelID == channelID && defaltChannelName == title && mention == "everyone"){
			return true
		}
	}
	return false
}