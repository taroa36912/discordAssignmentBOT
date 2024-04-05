package main

import (
	"fmt"
	"formbot/cmd"
	"formbot/cmd/delete"
	"formbot/cmd/form"
	"formbot/cmd/move"
	"formbot/cmd/zemi"
	"formbot/event"
	"formbot/function"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	EnvDiscordToken = "discord_token"
	EnvClientId     = "client_id"
)

var (
	mentions []*discordgo.ApplicationCommandOptionChoice
	channelChoices []*discordgo.ApplicationCommandOptionChoice
	categoryChoices []*discordgo.ApplicationCommandOptionChoice
)

func main() {
	loadEnv()
	discordToken := os.Getenv(EnvDiscordToken)
	clientId := os.Getenv(EnvClientId)
	var (
		Token   = "Bot " + discordToken
		BotName = "<@" + clientId + ">"
	)
	fmt.Println(Token)
	fmt.Println(BotName)

	discord, err := discordgo.New(Token)
	discord.Token = Token

	if err != nil {
		fmt.Println("failed to login")
		fmt.Println(err)
	}

	mentions = subfunc.GenerateMentionOptions(discord)
	channelChoices = subfunc.GenerateChannelOptions(discord)
	categoryChoices = subfunc.GenerateCategoryOptions(discord)

	//イベントハンドラを追加
	discord.AddHandler(event.CheckReminder)
	discord.AddHandler(event.CreateZemiMessage)
	discord.AddHandler(event.CheckZemiReaction)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	cmds := cmd.NewExec()
	formCmd := form.NewFormCmd()
	formCmd.SetOptions(mentions)
	cmds.Add(formCmd)
	deleteCmd := delete.NewDeleteCmd()
	cmds.Add(deleteCmd)
	zemiCmd := zemi.NewZemiCmd()
	cmds.Add(zemiCmd)
	moveCmd := move.NewMoveCmd()
	moveCmd.SetOptions(channelChoices, categoryChoices)
	cmds.Add(moveCmd)

	cmdHandler := cmds.Activate(discord)
	defer cmdHandler.Deactivate()

	// 毎日20時に関数を実行するためのタイマーを作成
	ExecuteAt20(event.CreateZemiMessage, discord)
	// 毎時0分に関数を実行するためのタイマーを作成
	ExecuteAt0(event.CheckReminder, discord)
	// 毎日8時に関数を実行するためのタイマーを作成
	ExecuteAt8(event.CheckZemiReaction, discord)

	//ここから終了コマンド
	defer discord.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	fmt.Println("Removing commands...")
	// コマンドを削除
	cmdHandler.Deactivate()

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf(".env can't be accepted: %v", err)
	}
	fmt.Println(".env was accepted")
}

// ExecuteAt20は、指定した関数を毎日20時に実行します。
func ExecuteAt20(f func(s *discordgo.Session, e *discordgo.Ready), s *discordgo.Session) {
	// 現在の時刻を取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}
	// 現在時刻を取得
	now := time.Now().In(location)
	// 次の日の20時の時刻を計算
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 20, 0, 0, 0, now.Location())

	// 次の日の20時までの時間を計算
	duration := next.Sub(now)

	// タイマーを設定して、次の日の20時に関数を実行
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		// 関数を実行
		f(s, nil)
		// 20時に再び関数を実行するために、次の日の20時までの時間を計算して再度実行
		ExecuteAt20(f, s)
	}()
}

// ExecuteAt8は、指定した関数を毎日8時に実行します。
func ExecuteAt8(f func(s *discordgo.Session, e *discordgo.Ready), s *discordgo.Session) {
	// 現在の時刻を取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}

	// 現在時刻を取得
	now := time.Now().In(location)
	// 次の日の8時の時刻を計算
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 8, 0, 0, 0, now.Location())

	// 次の日の8時までの時間を計算
	duration := next.Sub(now)

	// タイマーを設定して、次の日の8時に関数を実行
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		// 関数を実行
		f(s, nil)
		// 8時に再び関数を実行するために、次の日の8時までの時間を計算して再度実行
		ExecuteAt8(f, s)
	}()
}

// ExecuteAt0は、指定した関数を毎時0分に実行します。
func ExecuteAt0(f func(s *discordgo.Session, e *discordgo.Ready), s *discordgo.Session) {
	// 現在の時刻を取得
	// 日本標準時の場所情報を取得
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location:", err)
		return
	}

	// 現在時刻を取得
	now := time.Now().In(location)

	// 次の時間の0分の時刻を計算
	next := now.Truncate(time.Hour).Add(time.Hour)

	// 次の時間の0分までの時間を計算
	duration := next.Sub(now)

	// タイマーを設定して、次の時間の0分に関数を実行
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		// 関数を実行
		f(s, nil)
		// 0分に再び関数を実行するために、次の時間の0分までの時間を計算して再度実行
		ExecuteAt0(f, s)
	}()
}
