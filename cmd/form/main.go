package form

import (
	"formbot/cmd/form/subform"
	"formbot/function"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
	"fmt"
)

type FormCmd struct {
}

func NewFormCmd() FormCmd {
	return FormCmd{}
}

// コマンド作成&入力関数
func (n FormCmd) Info() *discordgo.ApplicationCommand {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
	}
	// 現在時刻を取得
	start := time.Now().In(location)
	end := start.AddDate(1, 0, 0) // 現在から1年後までの日時を生成
	options := subfunc.GenerateDateTimeOptions(start, end)
	return &discordgo.ApplicationCommand{
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
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "mention",
						Description: "メンション範囲",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "自分のみ", Value: "me"},
							{Name: "全員", Value: "everyone"},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add-once",
				Description: "課題期限通知時間の更新を行います.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "time",
						Description: "通知をする年月日時",
						Required:    true,
						Choices: options,
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
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "自分のみ", Value: "me"},
							{Name: "全員", Value: "everyone"},
						},
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