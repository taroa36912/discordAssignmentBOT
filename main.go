package main

import (
	"fmt"
	"formbot/cmd"
	"formbot/cmd/delete"
	"formbot/cmd/form"
	"formbot/cmd/nox"
	"formbot/event/checkForm"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	EnvDiscordToken = "discord_token"
	EnvClientId     = "client_id"
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

	//イベントハンドラを追加
	discord.AddHandler(checkform.CheckReminder)
	err = discord.Open()

	if err != nil {
		fmt.Println(err)
	}

	cmds := cmd.NewExec()
	formCmd := form.NewFormCmd()
	cmds.Add(formCmd)
	deleteCmd := delete.NewDeleteCmd()
	cmds.Add(deleteCmd)
	noxCmd := nox.NewNoxCmd()
	cmds.Add(noxCmd)

	cmdHandler := cmds.Activate(discord)
	defer cmdHandler.Deactivate()

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
