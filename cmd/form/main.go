package form

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"formbot/cmd/form/subform"
	"log"
	"os"
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
		view.HandleViewCommand(s, i, options)
	case "add":
		add.HandleAddCommand(s, i, options)
	case "update":
		update.HandleUpdateCommand(s, i, options)
	case "delete":
		delete.HandleDeleteCommand(s, i, options)
	}
}

func WriteToDataFile(channelID, channelName, time, day string) error {
	dataPath := os.Getenv(EnvDataPath)
	file, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open data file: %v", err)
	}
	defer file.Close()

	data := fmt.Sprintf("%s, %s, %s, %s\n", channelID, channelName, time, day)
	if _, err := file.WriteString(data); err != nil {
		return fmt.Errorf("failed to write to data file: %v", err)
	}

	return nil
}

func WeekEtoJ(day string) (string, error) {
	dayJ := ""
	switch day {
	case "Sunday":
		dayJ = "日"
	case "Monday":
		dayJ = "月"
	case "Tuesday":
		dayJ = "火"
	case "Wednesday":
		dayJ = "水"
	case "Thursday":
		dayJ = "木"
	case "Friday":
		dayJ = "金"
	case "Saturday":
		dayJ = "土"
	default:
		return "", fmt.Errorf("invalid day: %s", day)
	}
	return dayJ, nil
}
