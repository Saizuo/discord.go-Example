package discord_cmd

import (
    "fmt"
    "strings"
    "strconv"
    "encoding/json"

	fputils "fpbot/pkg/utils"
    // fpmodel "fpbot/pkg/model"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

const pointsChannelName = "bot-data"

var (
    adminPointsType string
    adminPointsUserID string
)

type adminPointsCommand struct {
	DiscordCommand
}

func (c *adminPointsCommand) run(args []string) {
    messageGuild, err := fputils.GetMessageGuild(c.Session, c.Message)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

	pointsChannel, err := fputils.GetChannelFromGuild(pointsChannelName, messageGuild)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
		return
	}

    userData, messageID, err := fputils.FindUserDataFromDiscordDataStore(c.Session, c.Message, pointsChannel.ID, adminPointsUserID)
    if err != nil {
        if len(messageID) > 0 {
            // Create new user data
            guildMember, err := c.Session.GuildMember(messageGuild.ID, adminPointsUserID)
            if err != nil {
                c.Session.ChannelMessageSend(
                    c.Message.ChannelID,
                    fmt.Sprintf("Unable to create and get user details for %s: %s", adminPointsUserID, err.Error()),
                )
                return
            }
            username := guildMember.User.Username
            if len(guildMember.Nick) > 0 {
                username = guildMember.Nick
            }
            userData = fputils.NewUserData(adminPointsUserID, username)
        } else {
            c.Session.ChannelMessageSend(
                c.Message.ChannelID,
                fmt.Sprintf("Unable get user details for %s: %s", adminPointsUserID, err.Error()),
            )
            return
        }
    }

    paramValue, err := strconv.Atoi(args[0])
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to convert param to number: %s", err.Error()))
    }

    switch strings.ToLower(adminPointsType) {
    case "add":
        userData.Points += paramValue
    case "remove", "rm":
        userData.Points += paramValue
    case "set":
        userData.Points = paramValue
    }

    jsonData, err := json.MarshalIndent(userData, "", "    ")
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Unable to marshal data: %s", err.Error()))
        return
    }

    c.Session.ChannelMessageDelete(pointsChannel.ID, messageID)
    c.Session.ChannelMessageSend(pointsChannel.ID, fmt.Sprint(string(jsonData)))

    c.Session.ChannelMessageSend(c.Message.ChannelID, "Probably a success")
}

func NewAdminPointsCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
	dc := &adminPointsCommand{
		DiscordCommand: DiscordCommand{
			Session: s,
			Message: m,
			BotData: b,
		},
	}

	c := &cobra.Command{
		Use:   "points <--type type> <--id user id> <points>",
		Short: "Modify user points",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dc.run(args)
		},
	}
    c.Flags().StringVarP(&adminPointsType, "type", "t", "", "Type of action to perform")
    c.Flags().StringVarP(&adminPointsUserID, "id", "i", "", "User ID to modify")
    c.MarkFlagRequired("type")
    c.MarkFlagRequired("id")

    modifyUsageFunc(c, s, m)

	return c
}
