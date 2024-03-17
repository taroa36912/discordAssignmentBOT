package subfunc

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"formbot/cmd/form/subform"
	"log"
	"os"
)

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
