package discord_cmd

import (
	"fmt"
	"strings"

	fputils "fpbot/pkg/utils"
    // fpmodel "fpbot/pkg/model"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

const (
	streamInfoChannelName      = "stream-info"
	streamInfoScheduleTemplate = `Streaming schedule for https://www.twitch.tv/team_youwin
All times in US Eastern time.
Streams may start/end later than listed.
--------------------------------------------------------
Sun: %s
Mon: %s
Tue: %s
Wed: %s
Thu: %s
Fri: %s
Sat: %s
--------------------------------------------------------`
)

type updateStreamInfoCommand struct {
	DiscordCommand
}

func (c *updateStreamInfoCommand) run(args []string) {
	if len(args) < 1 {
		c.Session.ChannelMessageSend(c.Message.ChannelID, "Not enough params specified.")
	}

	messageGuild, err := fputils.GetMessageGuild(c.Session, c.Message)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	streamInfoChannel, err := fputils.GetChannelFromGuild(streamInfoChannelName, messageGuild)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	streamInfoChannelMessages, err := c.Session.ChannelMessages(streamInfoChannel.ID, 1, "", "", "")
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	newSchedule := ""
	switch args[0] {
	case "default":
		newSchedule = fmt.Sprintf(streamInfoScheduleTemplate,
			"7-10pm (Programming)",
			"n/a",
			"7-10pm (Programming)",
			"9-11pm (Programming)",
			"7-10pm (Programming)",
			"n/a",
			"n/a",
		)
		if len(streamInfoChannelMessages) > 0 {
			c.Session.ChannelMessageDelete(streamInfoChannel.ID, streamInfoChannelMessages[0].ID)
		}

		c.Session.ChannelMessageSend(streamInfoChannel.ID, fmt.Sprintf("```%s```", newSchedule))
	default:
		if len(args) < 2 {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Not enough commands to update stream-info")
			return
		}

		if len(streamInfoChannelMessages) < 1 {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "stream-info channel must have a message to edit")
			return
		}

		lastMessageID := streamInfoChannelMessages[0].ID

		lastMessage, err := c.Session.ChannelMessage(streamInfoChannel.ID, lastMessageID)
		if err != nil {
			c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
			return
		}

		newSchedule, err = fputils.ReplaceStringAt(lastMessage.Content, fmt.Sprintf("%s ", args[0]), "\n", strings.Join(args[1:], " "))
		if err != nil {
			c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
			return
		}

		c.Session.ChannelMessageDelete(streamInfoChannel.ID, lastMessageID)
		c.Session.ChannelMessageSend(streamInfoChannel.ID, newSchedule)
	}

	c.Session.ChannelMessageSend(c.Message.ChannelID, "Updated stream-info")
}

func NewUpdateStreamInfoCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
	dc := &updateStreamInfoCommand{
		DiscordCommand: DiscordCommand{
			Session: s,
			Message: m,
			BotData: b,
		},
	}

	c := &cobra.Command{
		Use:   "update-stream-info <day> <description>",
		Short: "Update the stream info",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dc.run(args)
		},
	}
	
    modifyUsageFunc(c, s, m)

	return c
}
