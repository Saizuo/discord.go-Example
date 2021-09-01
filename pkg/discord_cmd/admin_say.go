package discord_cmd

import (
    "strings"

    fputils "fpbot/pkg/utils"
    // fpmodel "fpbot/pkg/model"

    "github.com/spf13/cobra"
    dgo "github.com/bwmarrin/discordgo"
)

type sayCommand struct {
    DiscordCommand
}

func (c *sayCommand) run(args []string) {
    if len(args) < 2 {
        c.Session.ChannelMessageSend(c.Message.ChannelID, "Not enough params specified. Need 2 params: Channel and Message")
    }

    messageGuild, err := fputils.GetMessageGuild(c.Session, c.Message)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
        return
    }

    channelToSendTo, err := fputils.GetChannelFromGuild(args[0], messageGuild)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
        return
    }

    c.Session.ChannelMessageSend(channelToSendTo.ID, strings.Join(args[1:], " "))
}

func NewSayCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
    dc := &sayCommand{
        DiscordCommand: DiscordCommand{
            Session: s,
            Message: m,
            BotData: b,
        },
    }

    c := &cobra.Command{
        Use: "say <channel> <some-text>",
        Short: "Have the bot say something",
        Args: cobra.MinimumNArgs(2),
        Run: func(cmd *cobra.Command, args []string){
            dc.run(args)
        },
    }
    
    modifyUsageFunc(c, s, m)

    return c
}
