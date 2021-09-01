package discord_cmd

import (
	"fmt"
	"strings"

	// "strconv"
	"encoding/json"

	fputils "fpbot/pkg/utils"
    // fpmodel "fpbot/pkg/model"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

const userDataChannelName = "bot-data"

type userDataCommand struct {
	DiscordCommand
}

func (c *userDataCommand) run(args []string) {
	messageGuild, err := fputils.GetMessageGuild(c.Session, c.Message)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	userDataChannel, err := fputils.GetChannelFromGuild(userDataChannelName, messageGuild)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	authorID := c.Message.Author.ID

	userData, messageID, err := fputils.FindUserDataFromDiscordDataStore(
		c.Session,
		c.Message,
		userDataChannel.ID,
		authorID,
	)
	if err != nil {
		if len(messageID) > 0 {
            if len(args) > 0 {
                if strings.EqualFold(args[0], "delete") {
                    // Abort early if we try to delete but there is no data
                    c.Session.ChannelMessageSend(c.Message.ChannelID, "No user data to delete")
                    return
                }
            }
			// Create new user data
			guildMember, err := c.Session.GuildMember(messageGuild.ID, authorID)
			if err != nil {
				c.Session.ChannelMessageSend(
					c.Message.ChannelID,
					fmt.Sprintf("Unable to create and get user details for %s: %s", authorID, err.Error()),
				)
				return
			}
			username := guildMember.User.Username
			if len(guildMember.Nick) > 0 {
				username = guildMember.Nick
			}
			userData = fputils.NewUserData(authorID, username)

			jsonData, err := json.MarshalIndent(userData, "", "    ")
			if err != nil {
				c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to marshal data: %s", err.Error()))
				return
			}
			c.Session.ChannelMessageSend(userDataChannel.ID, fmt.Sprint(string(jsonData)))
		} else {
			c.Session.ChannelMessageSend(
				c.Message.ChannelID,
				fmt.Sprintf("Unable get user details for %s: %s", authorID, err.Error()),
			)
			return
		}
	}

	switch len(args) {
	case 0:
		// Read data
		jsonData, err := json.MarshalIndent(userData, "", "    ")
		if err != nil {
			c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to marshal data: %s", err.Error()))
			return
		}

		c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("```%s```", string(jsonData)))
	case 2:
		// Modify data
		switch strings.ToLower(args[0]) {
		case "twitchusername", "twitch":
			userData.TwitchUsername = args[1]
			jsonData, err := json.MarshalIndent(userData, "", "    ")
			if err != nil {
				c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to marshal data: %s", err.Error()))
				return
			}

			c.Session.ChannelMessageDelete(userDataChannel.ID, messageID)
			c.Session.ChannelMessageSend(userDataChannel.ID, fmt.Sprint(string(jsonData)))

			c.Session.ChannelMessageSend(c.Message.ChannelID, "Probably a success")
		case "delete":
			if strings.EqualFold(args[1], "confirmed") {
				c.Session.ChannelMessageDelete(userDataChannel.ID, messageID)
			}
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Success")
		}
	default:
		c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unrecognized command"))
	}
}

func NewUserDataCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
	dc := &userDataCommand{
		DiscordCommand: DiscordCommand{
			Session: s,
			Message: m,
			BotData: b,
		},
	}

	c := &cobra.Command{
		Use:   "user-data ([twitch <username>] | [delete confirmed])",
		Short: "Interact with user data",
		Run: func(cmd *cobra.Command, args []string) {
			dc.run(args)
		},
	}

	modifyUsageFunc(c, s, m)

	return c
}
