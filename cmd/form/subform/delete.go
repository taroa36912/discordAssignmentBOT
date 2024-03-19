package sub

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleDeleteCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// サブコマンドが発動されたことを確認
	log.Println("deleteコマンドが発動されました。")

	// 選択肢を用意
	choices := generateReminderChoices()

	// 選択肢を含んだメッセージを送信
	msg := "削除する通知を選択してください："
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
		Components: [][]discordgo.MessageComponent{
			{
				&discordgo.ActionRow{
					Components: choices,
				},
			},
		},
	})
	if err != nil {
		log.Printf("failed to send follow-up message with choices, err: %v", err)
		return
	}
}

func generateReminderChoices() []*discordgo.MessageComponent {
	// 簡単な選択肢を用意する（a, b, c）
	choices := []*discordgo.MessageComponent{
		{
			Type: discordgo.ComponentTypeButton,
			Style: discordgo.ButtonPrimary,
			Label: "選択肢A",
			CustomID: "deleteChoiceA",
		},
		{
			Type: discordgo.ComponentTypeButton,
			Style: discordgo.ButtonPrimary,
			Label: "選択肢B",
			CustomID: "deleteChoiceB",
		},
		{
			Type: discordgo.ComponentTypeButton,
			Style: discordgo.ButtonPrimary,
			Label: "選択肢C",
			CustomID: "deleteChoiceC",
		},
	}
	return choices
}
