package zemi

import (
	"fmt"
	"log"
	"os"
	"time"
	"formbot/function"

	"github.com/bwmarrin/discordgo"
)

type ZemiCmd struct {
}

func NewZemiCmd() ZemiCmd {
	return ZemiCmd{}
}

var (
	zemiChannelID = os.Getenv("zemi_channel_id")
	zemiRoleID    = os.Getenv("zemi_role_id")
)

// コマンド作成&入力関数
func (n ZemiCmd) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "zemi",
		Description: "ゼミの日程が変更した際，こちらで再設定します. Botが生成したメッセージは削除してください.",
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
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "minute",
				Description: "ゼミ開始の分",
				Required:    true,
			},
		},
	}
}

func (n ZemiCmd) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// インタラクションがDMからのものであるかを確認します
	if i.Interaction.GuildID == "" {
		log.Printf("command invoked in a DM, ignoring...")
		return
	}
	options := i.ApplicationCommandData().Options

	// 回答が正しく得られなかった場合，終了
	if len(options) != 4 && len(options) != 5 {
		log.Printf("invalid options: %#v", options)
		return
	}

	// 処理を行っている間表示されるメッセージ
	followUp := discordgo.WebhookParams{
		Content: "ゼミ出席メッセージ追加中...",
		Flags:   discordgo.MessageFlagEphemeral,
	}
	followUpMsg, err := s.InteractionRespondWithFollowup(i.Interaction, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}

	year := options[0].IntValue()
	month := options[1].IntValue()
	day := options[2].IntValue()
	hour := options[3].IntValue()
	weekday := getWeekday(int(year), int(month), int(day))
	dayJ, err := subfunc.WeekEtoJ(weekday)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
		return
	}
	minute := int(options[4].IntValue())
	if minute < 0 || 59 < minute {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "分は0~59の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.InteractionResponseEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}

	// ゼミ出席メッセージの追加
	sentence := fmt.Sprintf("<@%s>```%d年%d月%d日%s曜日%d時%d分\n自主ゼミの出欠を取ります.\nリアクションをしてください.```", zemiRoleID, year, month, day, dayJ, hour, minute)
	msg, err := s.ChannelMessageSend(zemiChannelID, sentence)
	if err != nil {
		log.Println("Error sending message : ", err)
	}

	subfunc.WritetoFile("zemiMessage.txt", fmt.Sprintf("%s, %d, %d, %d, %s, %d, %d", msg.ID, year, month, day, weekday, hour, minute))

	// 通知追加中の表示を変更する
	finishFollowUpStr := "zemiコマンドが正しく発動されました."
	finishFollowUp := discordgo.WebhookEdit{
		Content: &finishFollowUpStr,
	}
	if _, err := s.InteractionResponseEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}

}

func getWeekday(year, month, day int) string {
	// 年月日からTime型の変数を作成
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// 曜日を取得
	weekday := date.Weekday()
	// 曜日を文字列に変換して返す
	return weekday.String()
}
