package form

import (
	"formbot/cmd/form/subform"
	"log"

	"github.com/bwmarrin/discordgo"
)

type FormCmd struct {
	mentions []*discordgo.ApplicationCommandOptionChoice
}

func NewFormCmd() *FormCmd {
	return &FormCmd{}
}

// optionsを設定するメソッド
func (f *FormCmd) SetOptions(mentions []*discordgo.ApplicationCommandOptionChoice) {
    f.mentions = mentions
}
// コマンド作成&入力関数
func (n *FormCmd) Info() *discordgo.ApplicationCommand {
	// 選択肢を使用してコマンドを作成
	command :=  &discordgo.ApplicationCommand{
		Name:        "form",
		Description: "このチャンネルの課題の締め切りを通知する設定を行います.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "課題期限通知時間の閲覧を行います.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add-weekly",
				Description: "課題期限毎週通知の追加を行います.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "hour",
						Description: "通知する時間",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "0時", Value: 0},
							{Name: "1時", Value: 1},
							{Name: "2時", Value: 2},
							{Name: "3時", Value: 3},
							{Name: "4時", Value: 4},
							{Name: "5時", Value: 5},
							{Name: "6時", Value: 6},
							{Name: "7時", Value: 7},
							{Name: "8時", Value: 8},
							{Name: "9時", Value: 9},
							{Name: "10時", Value: 10},
							{Name: "11時", Value: 11},
							{Name: "12時", Value: 12},
							{Name: "13時", Value: 13},
							{Name: "14時", Value: 14},
							{Name: "15時", Value: 15},
							{Name: "16時", Value: 16},
							{Name: "17時", Value: 17},
							{Name: "18時", Value: 18},
							{Name: "19時", Value: 19},
							{Name: "20時", Value: 20},
							{Name: "21時", Value: 21},
							{Name: "22時", Value: 22},
							{Name: "23時", Value: 23},
						},
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
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "mention",
						Description: "メンション範囲",
						Required:    true,
						Choices:     n.mentions,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add-once",
				Description: "課題期限通知時間の更新を行います.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "year",
						Description: "通知する年",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "2024年", Value: 2024},
							{Name: "2025年", Value: 2025},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "month",
						Description: "通知する月",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "1月", Value: 1},
							{Name: "2月", Value: 2},
							{Name: "3月", Value: 3},
							{Name: "4月", Value: 4},
							{Name: "5月", Value: 5},
							{Name: "6月", Value: 6},
							{Name: "7月", Value: 7},
							{Name: "8月", Value: 8},
							{Name: "9月", Value: 9},
							{Name: "10月", Value: 10},
							{Name: "11月", Value: 11},
							{Name: "12月", Value: 12},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "day",
						Description: "通知する日(1~31の間の整数を入力)",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "hour",
						Description: "通知する時間",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "0時", Value: 0},
							{Name: "1時", Value: 1},
							{Name: "2時", Value: 2},
							{Name: "3時", Value: 3},
							{Name: "4時", Value: 4},
							{Name: "5時", Value: 5},
							{Name: "6時", Value: 6},
							{Name: "7時", Value: 7},
							{Name: "8時", Value: 8},
							{Name: "9時", Value: 9},
							{Name: "10時", Value: 10},
							{Name: "11時", Value: 11},
							{Name: "12時", Value: 12},
							{Name: "13時", Value: 13},
							{Name: "14時", Value: 14},
							{Name: "15時", Value: 15},
							{Name: "16時", Value: 16},
							{Name: "17時", Value: 17},
							{Name: "18時", Value: 18},
							{Name: "19時", Value: 19},
							{Name: "20時", Value: 20},
							{Name: "21時", Value: 21},
							{Name: "22時", Value: 22},
							{Name: "23時", Value: 23},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "通知する内容を簡潔に書いてください",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "mention",
						Description: "メンション範囲",
						Required:    true,
						Choices:     n.mentions,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "課題期限通知時間の削除を行います.viewコマンドで通知番号を確認してください.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "delete-index",
						Description: "削除する通知番号",
						Required:    true,
					},
				},
			},
		},
	}
	return command
}

func (n FormCmd) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// インタラクションがDMからのものであるかを確認します
	if i.Interaction.GuildID == "" {
		log.Printf("command invoked in a DM, ignoring...")
		return
	}

	opts := i.ApplicationCommandData().Options

	// サブコマンドが正しく選択されていることを確認
	if len(opts) == 0 {
		log.Printf("invalid options: %#v", opts)
		return
	}

	// サブコマンドに応じて処理を振り分け
	subCommand := opts[0].Name
	options := opts[0].Options

	// 処理を始める前の手続き
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to do interaction response, err: %v", err)
		return
	}

	// それぞれのサブコマンド内での処理を行う
	switch subCommand {
	case "view", "add-weekly", "add-once", "delete":
		handleSubCommand(s, i, subCommand, options)
	default:
		log.Printf("invalid subcommand: %s", subCommand)
	}

}

func handleSubCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	subCommand string,
	options []*discordgo.ApplicationCommandInteractionDataOption,
) {
	// サブコマンドに応じて処理を振り分け
	switch subCommand {
	case "view":
		sub.HandleViewCommand(s, i, options)
	case "add-weekly":
		sub.HandleAddWeeklyCommand(s, i, options)
	case "add-once":
		sub.HandleAddOnceCommand(s, i, options)
	case "delete":
		sub.HandleDeleteCommand(s, i, options)
	}
}
