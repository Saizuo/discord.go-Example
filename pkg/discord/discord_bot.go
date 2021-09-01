package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
    "time"

    // fpmodel "fpbot/pkg/model"
    cmd "fpbot/pkg/discord_cmd"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	discordToken := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Error creating Discord session, ", err)
		return
	}

	defer dg.Close()

	bd := cmd.BotData{
		StartTime:    time.Now(),
		LastRateLimitedCommandTime: time.Now(),
	}
    // bd := NewBotData()
	dg.AddHandler(bd.HandleRegularText)

    as := NewAntiSpam()
    dg.AddHandler(as.handleSpam)

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection, ", err)
		return
	}

	fmt.Println("Bot is running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
