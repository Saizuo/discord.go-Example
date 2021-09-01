package discord_cmd

import (
	"fmt"

    // fpmodel "fpbot/pkg/model"

	dgo "github.com/bwmarrin/discordgo"
    db "github.com/replit/database-go"
	"github.com/spf13/cobra"
)

type dbDeleteCommand struct {
	DiscordCommand
}

func (c *dbDeleteCommand) run(args []string) {
    keyToDelete := args[0]

	_, err := db.Get(keyToDelete)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("%s does not exist in db: %s", keyToDelete, err.Error()))
        return
    }

    err = db.Delete(keyToDelete)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to delete key %s: %s", keyToDelete, err.Error()))
        return
    }

    c.Session.ChannelMessageSend(c.Message.ChannelID, "Probably a success")
}

func NewDBDeleteCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
	dc := &dbDeleteCommand{
		DiscordCommand: DiscordCommand{
			Session: s,
			Message: m,
			BotData: b,
		},
	}

	c := &cobra.Command{
		Use:   "db-delete <key>",
		Short: "Deletes an item from the DB",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dc.run(args)
		},
	}
	
    modifyUsageFunc(c, s, m)

	return c
}
