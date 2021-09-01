package twitch

import (
	"log"

	twitchgo "github.com/gempir/go-twitch-irc/v2"
)

func Run() {
	client := twitchgo.NewAnonymousClient()

	log.Println(client)
}
