package form

import (
	"fmt"
	"formbot/send/send_e"
	"github.com/bwmarrin/discordgo"
	"log"
)

type FormCmd struct {
}

func NewFormCmd() FormCmd {
	return FormCmd{}
}

func (n FormCmd) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "form",
		Description: "このチャンネルに，このチャンネルの課題の締め切りを通知する設定を行います.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "hour",
				Description: "通知する時間(時)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "day",
				Description: "通知する曜日",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "日曜日", Value: "Sunday"},
					{Name: "月曜日", Value: "Monday"},
					{Name: "火曜日", Value: "Tuesday"},
					{Name: "水曜日", Value: "Wednesday"},
					{Name: "木曜日", Value: "Thursday"},
					{Name: "金曜日", Value: "Friday"},
					{Name: "土曜日", Value: "Saturday"},
				},
			},
		},
	}
}

func (n FormCmd) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	opts := i.ApplicationCommandData().Options
	if len(opts) != 2 {
		log.Printf("invalid options: %#v", opts)
		return
	}
	count := opts[0].IntValue()
	day := opts[1].StringValue()
	if count < 0 || 23 < count {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "時間が無効です. 0~23の整数で選択してください",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Printf("failed to respond, err: %v", err)
		}
		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to do interaction response, err: %v", err)
		return
	}

	// create follow-up message
	followUp := discordgo.WebhookParams{
		Content: "送信中...",
		Flags:   discordgo.MessageFlagsEphemeral,
	}

	followUpMsg, err := s.FollowupMessageCreate(i.Interaction, true, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}

	// dayオブジェクトを使った処理の関数
	result := fmt.Sprintf("%d : %s", count, day)
	send.SendMessage(s, i.ChannelID, result)


	finishFollowUpStr := "メッセージが送信されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}

	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}

	// Now you can proceed with the rest of your code, utilizing the 'count' and 'day' variables.
}
