package subfunc

import (
	//"github.com/bwmarrin/discordgo"
	//"log"
	"os"
	"fmt"
	"bufio"
)


const (
	EnvZemiChannel = "zemi_channel_id"
)


// 指定したファイル名に，指定した文章を書き込み
func WriteToDataFile(fileName string, str string) error {
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


