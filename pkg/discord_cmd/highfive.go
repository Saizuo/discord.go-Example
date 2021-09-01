package discord_cmd

import (
    "strconv"
    "strings"

    // fpmodel "fpbot/pkg/model"

    "github.com/spf13/cobra"
    dgo "github.com/bwmarrin/discordgo"
)

type highFiveCommand struct {
    DiscordCommand
}

func (c *highFiveCommand) run() {
    emoji, err := strconv.ParseInt(strings.TrimPrefix("\\U0001f44f", "\\U"), 16, 32)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
        return
    }
    c.Session.ChannelMessageSend(c.Message.ChannelID, string(emoji))
}

func NewHighFiveCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
    dc := &highFiveCommand{
        DiscordCommand: DiscordCommand{
            Session: s,
            Message: m,
            BotData: b,
        },
    }

    c := &cobra.Command{
        Use: "high-five",
        Short: "High five the bot",
        Run: func(cmd *cobra.Command, args []string){
            dc.run()
        },
    }
    
    modifyUsageFunc(c, s, m)

    return c
}
