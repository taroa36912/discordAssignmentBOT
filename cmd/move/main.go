package move

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type MoveCmd struct {
	channelChoices []*discordgo.ApplicationCommandOptionChoice
	categoryChoices []*discordgo.ApplicationCommandOptionChoice
}

func NewMoveCmd() *MoveCmd {
	return &MoveCmd{}
}

// optionsを設定するメソッド
func (f *MoveCmd) SetOptions(channelChoices []*discordgo.ApplicationCommandOptionChoice, categoryChoices []*discordgo.ApplicationCommandOptionChoice) {
    f.channelChoices = channelChoices
	f.categoryChoices = categoryChoices
}

func (n *MoveCmd) Info() *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "channel",
		Description: "移動するチャンネルを選択してください.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "channel",
				Description: "移動するチャンネル",
				Required:    true,
				Choices: n.channelChoices,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "category",
				Description: "移動先のカテゴリー",
				Required:    true,
				Choices: n.categoryChoices,
			},
		},
	}
	return command
}

func (n MoveCmd) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	opts := i.ApplicationCommandData().Options
	if len(opts) != 2 {
		log.Printf("invalid options: %#v", opts)
		return
	}
	channelID := opts[0].StringValue()
	categoryID := opts[1].StringValue()


	// create follow-up message
	followUp := discordgo.WebhookParams{
		Content: "移動中...",
		Flags:   discordgo.MessageFlagsEphemeral,
	}
	followUpMsg, err := s.FollowupMessageCreate(i.Interaction, true, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}



	// Move the channel to the specified category
	_, err = s.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
    	ParentID: categoryID,
	})
	if err != nil {
    	s.ChannelMessageSend(i.ChannelID, "Error moving channel: "+err.Error())
    	return
	}


	finishFollowUpStr := "チャンネルを移動しました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}
}
