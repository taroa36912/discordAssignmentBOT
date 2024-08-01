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
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "ゼミ内容を簡潔に書いてください",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "place",
				Description: "ゼミの開催場所を書いてください",
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
	if len(options) != 7 {
		log.Printf("invalid options: %#v", options)
		return
	}

	zemiChannelID := os.Getenv("zemi_channel_id")
	zemiRoleID    := os.Getenv("zemi_role_id")
	zemiAttendID := os.Getenv("emoji_attend_id")
	emojiAttend := fmt.Sprintf(":shusseki:%s", zemiAttendID)
	zemiAbsentID := os.Getenv("emoji_absent_id")
	emojiAbsent := fmt.Sprintf(":kesseki:%s", zemiAbsentID)

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
		// エラーメッセージを送信します
		errMsg := "分は0~59の範囲で入力してください."
		sendError(s, i, errMsg)
		return
	}

	message := options[5].StringValue()

	place := options[6].StringValue()

	// ゼミ出席メッセージの追加
	sentence := fmt.Sprintf("<@&%s>```自主ゼミ開催通知\n日時 : %d年%d月%d日%s曜日%d時%d分\n内容 : %s, 場所 : %s\n```", zemiRoleID, year, month, day, dayJ, hour, minute, message, place)
	mes, err := s.ChannelMessageSend(zemiChannelID, sentence)
	if err != nil {
		log.Println("Error sending message : ", err, zemiChannelID, sentence)
		return
	}

	err = s.MessageReactionAdd(zemiChannelID, mes.ID, emojiAbsent)
	if err != nil {
		log.Println("Error adding reaction,", err)
		return
	}

	err = s.MessageReactionAdd(zemiChannelID, mes.ID, emojiAttend)
	if err != nil {
		log.Println("Error adding reaction,", err)
		return
	}

	subfunc.WritetoFile("zemiMessage.txt", fmt.Sprintf("%s, %d, %d, %d, %s, %d, %d, %s, %s", mes.ID, year,  month, day, weekday, hour, minute, message, place))

	// 正常なリクエストの返信
	successMsg := "zemiコマンドが正しく発動されました."
	sendEphemeralResponse(s, i, successMsg)
}

func getWeekday(year, month, day int) string {
	// 年月日からTime型の変数を作成
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// 曜日を取得
	weekday := date.Weekday()
	// 曜日を文字列に変換して返す
	return weekday.String()
}

// エラーメッセージを送信するヘルパー関数
func sendError(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errMsg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("failed to send error message: %v", err)
	}
}

// 正常なリクエストの返信を送信するヘルパー関数
func sendEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("failed to send response message: %v", err)
	}
}
