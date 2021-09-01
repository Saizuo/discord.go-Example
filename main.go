package main

import (
	"fmt"
    "os/signal"
    "os"
    "syscall"

	"fpbot/pkg/discord"
)

func main() {
	fmt.Println("Hello world")
	go discord.Run()

    sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
