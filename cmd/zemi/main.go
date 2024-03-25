package zemi

import (
	"fmt"
	subfunc "formbot/function"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ZemiCmd struct {
}

func NewZemiCmd() ZemiCmd {
	return ZemiCmd{}
}

const (
	EnvZemiChannel = "zemi_channel_id"
)

// コマンド作成&入力関数
func (n ZemiCmd) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "zemi",
		Description: "ゼミの日程・曜日が変更した際，こちらで再設定します. ゼミの先生, 年月日を入力してください. Botが生成したメッセージは削除してください.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ゼミの先生(敬称略)",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "野崎", Value: "1221179383742595204"},
					{Name: "その他", Value: "1221339641031032842"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "year",
				Description: "ゼミの年",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "month",
				Description: "ゼミの月",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "day",
				Description: "ゼミの日",
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
	if len(options) != 4 {
		log.Printf("invalid options: %#v", options)
		return
	}

	zemiChannelID := os.Getenv(EnvZemiChannel)
	mentionID := options[0].StringValue()
	year := options[1].IntValue()
	month := options[2].IntValue()
	day := options[3].IntValue()
	teacher, err := subfunc.GetTeacherName(mentionID)
	if err != nil {
		log.Printf("failed to get teacher name: %v", err)
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

	// 時間が正しく入力されなかった場合，終了
	if year < 2024 {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "年号2024~の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}
	if month < 1 || 12 < month {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "月は1~12の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}
	if day < 1 || 31 < day {
		// 通知追加中の表示を変更する
		finishFollowUpStr := "日は1~31の範囲で入力してください."
		finishFollowUp := discordgo.WebhookEdit{
			Content: &finishFollowUpStr,
		}
		if _, err := s.FollowupMessageEdit(i.Interaction, followUpMsg.ID, &finishFollowUp); err != nil {
			log.Printf("failed to edit follow-up message, err: %v", err)
			return
		}
		return
	}

	// 曜日の取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, location)
	weekday := date.Weekday().String()
	dayJ, err := subfunc.WeekEtoJ(weekday)
	if err != nil {
		log.Printf("failed to convert day to Japanese: %v", err)
		return
	}

	// ゼミ出席メッセージの追加
	sentence := fmt.Sprintf("<@%s>```%d年%d月%d日%s曜日\n%sゼミの出欠を取ります.\nリアクションをしてください.```", mentionID, year, month, day, dayJ, teacher)
	msg, err := s.ChannelMessageSend(zemiChannelID, sentence)
	if err != nil {
		log.Println("Error sending message : ", err)
	}

	subfunc.WriteFile("zemiMessage.txt", fmt.Sprintf("%s, %s, %d, %d, %d, %s", msg.ID, mentionID, year, month, day, weekday))

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
