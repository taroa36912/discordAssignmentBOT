package subfunc

import (
	"fmt"
	"os"
)

func WriteToDataFile(channelID, channelName, time, day string) error {
	// OpenFile関数でファイルを開きます。第二引数にos.O_APPEND|os.O_WRONLY|os.O_CREATEを指定することで、
	// 書き込み専用で、ファイルが存在しない場合は新しく作成し、ファイルの末尾に追記するように設定しています。
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open data file: %v", err)
	}
	defer file.Close()

	// 書き込むデータを作成します
	data := fmt.Sprintf("%s, %s, %s, %s\n", channelID, channelName, time, day)

	// WriteStringメソッドを使ってファイルにデータを書き込みます
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
