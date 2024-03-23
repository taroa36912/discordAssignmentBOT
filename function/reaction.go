package subfunc

import (
	//"github.com/bwmarrin/discordgo"
	//"log"
	"os"
	"fmt"
	"bufio"
)

func ReadZemiFile() ([]string, error) {
	// ファイルを開きます。ファイルが存在しない場合はエラーを返します。
	file, err := os.Open("zemi.txt")
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