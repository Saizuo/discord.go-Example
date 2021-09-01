package discord_cmd

import (
    // fpmodel "fpbot/pkg/model"

    "github.com/spf13/cobra"
    dgo "github.com/bwmarrin/discordgo"
)

type whoamiCommand struct {
    DiscordCommand
}

func (c *whoamiCommand) run() {
    whoamiName := c.Message.Author.Username
    if len(c.Message.Member.Nick) > 0 {
        whoamiName = c.Message.Member.Nick
    }
    c.Session.ChannelMessageSend(c.Message.ChannelID, whoamiName)
}

func NewWhoAmICommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
    dc := &whoamiCommand{
        DiscordCommand: DiscordCommand{
            Session: s,
            Message: m,
            BotData: b,
        },
    }

    c := &cobra.Command{
        Use: "whoami",
        Short: "Find out who you are",
        // Long: "Finds information about the current user",
        Args: cobra.ExactArgs(0),
        Run: func(cmd *cobra.Command, args []string){
            dc.run()
        },
    }
    c.SetOut(dc)
    c.SetErr(dc)

    return c
}
