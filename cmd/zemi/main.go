package zemi

import (
	"fmt"
	"formbot/function"
	"log"
	"os"
	"time"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type ZemiCmd struct {
}

func NewZemiCmd() ZemiCmd {
	return ZemiCmd{}
}

var (
	zemiChannelID = os.Getenv("myserver_zemi_channel_id")
	zemiRoleID    = os.Getenv("myserver_zemi_role_id")
)

// コマンド作成&入力関数
func (n ZemiCmd) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "zemi",
		Description: "ゼミの日程が変更した際，こちらで再設定します. Botが生成したメッセージは削除してください.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "ゼミ開始の年月日時",
				Required:    true,
				Choices:     subfunc.GenerateDateTimeOptions(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "minute",
				Description: "ゼミ開始の分",
				Required:    false,
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
	if len(options) != 2 {
		log.Printf("invalid options: %#v", options)
		return
	}

	// 処理を行っている間表示されるメッセージ
	followUp := discordgo.WebhookParams{
		Content: "ゼミ出席メッセージ追加中...",
		Flags:   discordgo.MessageFlagsEphemeral,
	}
	followUpMsg, err := s.FollowupMessageCreate(i.Interaction, true, &followUp)
	if err != nil {
		log.Printf("failed to send follow-up message, err: %v", err)
		return
	}

	dataString := options[0].StringValue()
	minute := int(options[1].IntValue())
	minuteStr := strconv.Itoa(minute)
	if minuteStr == ""{
		minute = 0
	}else if minute < 0 || 59 < minute {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "分は0~59の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}

	// time.Parse を使って時間を解析し、変数に格納
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dataString)
	if err != nil {
		fmt.Println("Failed to parse time:", err)
		return
	}

	// 年、月、日、時を取得して変数に格納
	year := parsedTime.Year()
	month := parsedTime.Month()
	day := parsedTime.Day()
	hour := parsedTime.Hour()

	// 指定された形式の文字列をパース
    t, err := time.Parse("2006-01-02T15:04:05", dataString)
    if err != nil {
        fmt.Println("Failed to parse time:", err)
        return
    }
    
    // 曜日を取得
    weekday := t.Weekday().String()
	dayJ, err := subfunc.WeekEtoJ(weekday)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
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
	if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
		log.Printf("failed to edit follow-up message, err: %v", err)
		return
	}

}
