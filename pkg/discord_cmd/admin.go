package discord_cmd

import (
    fputils "fpbot/pkg/utils"
    // fpmodel "fpbot/pkg/model"

    "github.com/spf13/cobra"
    dgo "github.com/bwmarrin/discordgo"
)

type adminCommand struct {
    DiscordCommand
}

func NewAdminCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
    c := &cobra.Command{
        Use: "admin",
        Short: "Perform an admin action",
        Hidden: true,
        Args: cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string){
            validRole := fputils.CheckForRole("Admin", s, m)

            if !validRole {
                s.ChannelMessageSend(m.ChannelID, "Invalid role, aborting.")
                return
            }
        },
    }
    
    c.AddCommand(
        NewSayCommand(s, m, b),
        NewUpdateStreamInfoCommand(s, m, b),
        NewDBDeleteCommand(s, m, b),
        NewAdminPointsCommand(s, m, b),
    )

    modifyUsageFunc(c, s, m)

    return c
}
