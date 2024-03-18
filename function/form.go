package subfunc

import (
	"log"
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"github.com/bwmarrin/discordgo"
)

func WriteToDataFile(str string) error {
	// OpenFile関数でファイルを開きます。第二引数にos.O_APPEND|os.O_WRONLY|os.O_CREATEを指定することで、
	// 書き込み専用で、ファイルが存在しない場合は新しく作成し、ファイルの末尾に追記するように設定しています。
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

func ReadDataFile() ([]string, error) {
	// ファイルを開きます。ファイルが存在しない場合はエラーを返します。
	file, err := os.Open("data.txt")
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


func CheckEachRow(s *discordgo.Session, e *discordgo.Ready, data string, current time.Time)(string, string){
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ5はweekly
	if len(parts) == 5 {
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		hour := parts[2]
		weekday := parts[3]
		mention := parts[4]

		currentHour := current.Hour()
		currentWeekday := current.Weekday().String()

		// 記録された時刻と曜日と現在の時刻と曜日が一致する場合にSendMessageを実行
		if hour == fmt.Sprintf("%02d", currentHour) && weekday == currentWeekday {
			if(mention == "everyone"){
				sentence := fmt.Sprintf("@everyone\n```課題：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			}else if(mention == "me"){
				sentence := fmt.Sprintf("```課題：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			}
		}
	}else if len(parts) == 7 { // 長さ7はonce
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		year := parts[2]
		month := parts[3]
		day := parts[4]
		hour := parts[5]
		mention := parts[6]

		currentHour := current.Hour()
		currentYear := current.Year()
		currentMonth := int(current.Month())
		currentDay := current.Day()

		// 指定された年月日時が現在時刻と一致する場合にのみ処理を実行
		if year == fmt.Sprintf("%d", currentYear) &&
			month == fmt.Sprintf("%d", currentMonth) &&
			day == fmt.Sprintf("%d", currentDay) &&
			hour == fmt.Sprintf("%d", currentHour) {
			if mention == "everyone" {
				sentence := fmt.Sprintf("@everyone\n```課題：%sの締め切りを通知します.\n```", title)
				return channelID, sentence
			} else if mention == "me" {
				sentence := fmt.Sprintf("課題：%sの締め切りを通知します.", title)
				return channelID, sentence
			}
		}
	}
	return "", ""
}


func ViewEachRow(myChannelID string, data string)(string){
	// データを", "で分割
	parts := strings.Split(data, ", ")

	// データの長さで, weeklyか, onceかを判別
	// 長さ5はweekly
	if len(parts) == 5 {
		// 各要素を変数に格納
		ChannelID := parts[0]
		title := parts[1]
		hour := parts[2]
		weekday := parts[3]
		mention := parts[4]
		dayJ, err := WeekEtoJ(weekday)
		if err != nil {
			log.Printf("failed to convert day to Japanese: %v", err)
			return ""
		}
		// 自分がメンション対象の時，表示
		if(mention == "everyone"){
			sentence := fmt.Sprintf("```形式 : 毎週(weekly)\nメンション対象 : %s\n課題 : %s\n%s曜日 : %s時\n```",mention, title, dayJ, hour)
			return sentence
		}else if(mention == "me"){
			if(myChannelID == ChannelID){
				sentence := fmt.Sprintf("```形式 : 毎週(weekly)\nメンション対象 : %s\n課題 : %s\n%s曜日 : %s時\n```",mention, title, dayJ, hour)
				return sentence
			}
		}
	}else if len(parts) == 7 { // 長さ7はonce
		// 各要素を変数に格納
		channelID := parts[0]
		title := parts[1]
		year := parts[2]
		month := parts[3]
		day := parts[4]
		hour := parts[5]
		mention := parts[6]

		// 指定された年月日時が現在時刻と一致する場合にのみ処理を実行
		if mention == "everyone" {
			sentence := fmt.Sprintf("```形式 : 一回のみ(once)\nメンション対象 : %s\n課題 : %s\n%s年%s月%s日%s時\n```",mention, title, year, month, day, hour)
			return sentence
		} else if mention == "me" {
			if(myChannelID == channelID){
				sentence := fmt.Sprintf("```形式 : 一回のみ(once)\nメンション対象 : %s\n課題 : %s\n%s年%s月%s日%s時\n```",mention, title, year, month, day, hour)
				return sentence
			}
		}
	}
	return  ""

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
