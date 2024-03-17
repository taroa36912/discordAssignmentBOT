package form

import (
	"github.com/bwmarrin/discordgo"
	"formbot/cmd/form/subform"
	"log"
)

const EnvDataPath = "data_path"
type FormCmd struct {
}

func NewFormCmd() FormCmd {
	return FormCmd{}
}

// コマンド作成&入力関数
func (n FormCmd) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "form",
		Description: "このチャンネルの課題の締め切りを通知する設定を行います.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "課題期限通知時間の閲覧を行います.",
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
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "課題期限通知時間の追加を行います.",
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
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "update",
				Description: "課題期限通知時間の更新を行います.",
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
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "課題期限通知時間の削除を行います.",
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
			},
		},
	}
}


func (n FormCmd) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	opts := i.ApplicationCommandData().Options

	// サブコマンドが正しく選択されていることを確認
	if len(opts) == 0 {
		log.Printf("invalid options: %#v", opts)
		return
	}

	// サブコマンドに応じて処理を振り分け
	subCommand := opts[0].Name
	switch subCommand {
	case "view", "add", "update", "delete":
		handleSubCommand(s, i, subCommand, opts[0].Options)
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
	case "add":
		sub.HandleAddCommand(s, i, options)
	case "update":
		sub.HandleUpdateCommand(s, i, options)
	case "delete":
		sub.HandleDeleteCommand(s, i, options)
	}
}