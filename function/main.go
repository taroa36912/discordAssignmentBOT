package subfunc

import (
	"github.com/bwmarrin/discordgo"
	//"log"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	EnvZemiChannel = "zemi_channel_id"
)

// 指定したファイルの，strに一致するものを削除する
func DeleteFile(fileName string, str string) error {
	// ファイルを開きます。ファイルが存在しない場合はエラーを返します。
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open data file: %v", err)
	}
	defer file.Close()

	var newLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 特定の文字列と一致する行を削除します
		if line != str {
			newLines = append(newLines, line)
		}
	}

	// スキャン中にエラーが発生した場合はエラーを返します。
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan data file: %v", err)
	}

	// ファイルを閉じてから、新しい内容を書き込みます
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close data file: %v", err)
	}

	// ファイルを再度開いて、新しい内容を書き込みます
	if err := os.WriteFile(fileName, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write to data file: %v", err)
	}

	return nil
}

// 自分のメンションのみtrueを返す
func ViewEachRow(myChannelID string, data string) string {
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ6はweekly
	if len(parts) == 6 {
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		hour := parts[2]
		weekday := parts[3]
		mention := parts[4]
		mentionName := parts[5]
		dayJ, err := WeekEtoJ(weekday)
		if err != nil {
			log.Printf("failed to convert day to Japanese: %v", err)
			return ""
		}
		// 自分個人のメンション，またはロールへのメンションだった場合，返信
		if mention == "me" {
			if myChannelID == channelID {
				sentence := fmt.Sprintf("形式 : 毎週(weekly)\nメンション対象 : 自分のみ\n課題 : %s\n%s曜日 : %s時\n", title, dayJ, hour)
				return sentence
			} else { // 他人個人のメンションは表示しない
				return ""
			}
		} else {
			sentence := fmt.Sprintf("形式 : 毎週(weekly)\nメンション対象 : %s\n課題 : %s\n%s曜日 : %s時\n", mentionName, title, dayJ, hour)
			return sentence
		}
	} else if len(parts) == 8 { // 長さ8はonce
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		year := parts[2]
		month := parts[3]
		day := parts[4]
		hour := parts[5]
		mention := parts[6]
		mentionName := parts[7]

		// 自分個人のメンション，またはロールへのメンションだった場合，返信
		if mention == "me" {
			if myChannelID == channelID {
				sentence := fmt.Sprintf("形式 : 一回のみ(once)\nメンション対象 : %s\n課題 : %s\n%s年%s月%s日%s時\n", mentionName, title, year, month, day, hour)
				return sentence
			} else { // 他人個人のメンションは表示しない
				return ""
			}
		} else {
			sentence := fmt.Sprintf("形式 : 一回のみ(once)\nメンション対象 : %s\n課題 : %s\n%s年%s月%s日%s時\n", mentionName, title, year, month, day, hour)
			return sentence
		}
	}
	return ""

}

func MentionType(data string) string {
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ6はweekly
	if len(parts) == 6 {
		mention := parts[4]
		return mention
	} else if len(parts) == 8 { // 長さ8はonce
		mention := parts[6]
		return mention
	}
	return ""
}

// ロールIDからロール名を取得する関数
func GetRoleName(s *discordgo.Session, guildID, roleID string) (string, error) {
	// サーバーからロール情報を取得
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "", err
	}

	// ロールIDに対応するロール名を探す
	for _, role := range roles {
		if role.ID == roleID {
			return role.Name, nil
		}
	}

	// ロールIDに対応するロールが見つからなかった場合
	return "", fmt.Errorf("ロールID %s に対応するロールが見つかりませんでした", roleID)
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

// 指定したファイル名に，指定した文章を書き込み
func WritetoFile(fileName string, str string) error {
	// OpenFile関数でファイルを開きます。第二引数にos.O_APPEND|os.O_WRONLY|os.O_CREATEを指定することで、
	// 書き込み専用で、ファイルが存在しない場合は新しく作成し、ファイルの末尾に追記するように設定しています。
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open data file: %v", err)
	}
	defer file.Close()

	// ファイルの最終行を読み込む
	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	// 最終行の下に書き込むデータを作成します
	sentence := fmt.Sprintf("\n%s", str)

	// WriteStringメソッドを使ってファイルにデータを書き込みます
	if _, err := file.WriteString(sentence); err != nil {
		return fmt.Errorf("failed to write to data file: %v", err)
	}

	// 最終行が空でない場合、改行を追加してから新しいデータを追記します
	if lastLine != "" {
		if _, err := file.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline to data file: %v", err)
		}
	}

	return nil
}

// 指定したファイル名のファイルを取得
func ReadFile(fileName string) ([]string, error) {
	// ファイルを開きます。ファイルが存在しない場合はエラーを返します。
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// スキャン中にエラーが発生した場合はエラーを返します。
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan data file: %v", err)
	}

	return lines, nil
}

// サーバーからロール情報を取得し，選択肢とする関数
func GenerateMentionOptions(s *discordgo.Session) []*discordgo.ApplicationCommandOptionChoice {
	// サーバーID
	serverID := os.Getenv("guild_id")
	var mentions []*discordgo.ApplicationCommandOptionChoice

	// サーバーからロール情報を取得
	roles, err := s.GuildRoles(serverID)
	if err != nil {
		log.Fatal("Error getting guild roles: ", err)
		return mentions
	}

	// ロール情報を選択肢に変換
	for _, role := range roles {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  role.Name,
			Value: role.ID,
		}
		mentions = append(mentions, choice)
		if len(mentions) > 24 {
			break
		}
	}
	choice := &discordgo.ApplicationCommandOptionChoice{
		Name:  "自分のみ",
		Value: "me",
	}
	mentions = append(mentions, choice)
	return mentions
}

// サーバーからチャンネル情報を取得し，選択肢とする関数
func GenerateChannelOptions(s *discordgo.Session) []*discordgo.ApplicationCommandOptionChoice {
	// サーバーID
	serverID := os.Getenv("guild_id")

	// チャンネル一覧を取得
	channels, err := s.GuildChannels(serverID)
	if err != nil {
		log.Fatal("Error getting guild channels: ", err)
	}

	// チャンネル一覧をコマンドの選択肢に変換
	var channelChoices []*discordgo.ApplicationCommandOptionChoice
	for _, channel := range channels {
		// テキストチャンネルのみを選択肢とする場合
		if channel.Type == discordgo.ChannelTypeGuildText {
			choice := &discordgo.ApplicationCommandOptionChoice{
				Name:  channel.Name,
				Value: channel.ID,
			}
			channelChoices = append(channelChoices, choice)
			if len(channelChoices) > 24 {
				break
			}
		}
	}

	return channelChoices
}

// サーバーからカテゴリー情報を取得し，選択肢とする関数
func GenerateCategoryOptions(s *discordgo.Session) []*discordgo.ApplicationCommandOptionChoice {
	// サーバーID
	serverID := os.Getenv("guild_id")
	// カテゴリー一覧を取得
	channels, err := s.GuildChannels(serverID)
	if err != nil {
		log.Fatal("Error getting guild channels: ", err)
	}

	// カテゴリー一覧をコマンドの選択肢に変換
	var categoryChoices []*discordgo.ApplicationCommandOptionChoice
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			choice := &discordgo.ApplicationCommandOptionChoice{
				Name:  channel.Name,
				Value: channel.ID,
			}
			categoryChoices = append(categoryChoices, choice)
			if len(categoryChoices) > 24 {
				break
			}
		}
	}

	return categoryChoices
}
