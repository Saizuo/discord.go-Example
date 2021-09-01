package utils

import (
    "strings"
    "errors"
    "fmt"
    "encoding/json"

    dgo "github.com/bwmarrin/discordgo"
)

const (
    DiscordBotPrefix = "?"
)

type UserData struct {
    UserID string
    Username string
    TwitchUsername string
    Points int
    Warnings uint
    Notes string
}

func NewUserData(userID string, username string) *UserData {
    return &UserData {
        UserID: userID,
        Username: username,
        TwitchUsername: "",
        Points: 0,
        Warnings: 0,
        Notes: "",
    }
}

// FindUserDataFromDiscordDataStore : Returns a user data struct and the discord message id
func FindUserDataFromDiscordDataStore(s *dgo.Session, m *dgo.Message, dataChannelID string, userID string) (*UserData, string, error) {
    lastMessageInChannelResults, err := s.ChannelMessages(dataChannelID, 1, "", "", "")
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Unable to get last message in user data channel")
        return nil, "", err
    }

    if len(lastMessageInChannelResults) < 1 {
        noDataErrorText := "No data in user data channel"
        s.ChannelMessageSend(m.ChannelID, noDataErrorText)
        return nil, "n/a", errors.New(noDataErrorText)
    }

    lastMessageInChannel := lastMessageInChannelResults[0]

    var userData *UserData

    userData, err = UserDataFromDiscordDataStoreMessage(lastMessageInChannel)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, err.Error())
        return nil, "", err
    }

    if userData.UserID == userID {
        return userData, lastMessageInChannel.ID, nil
    }

    beforeMessageID := lastMessageInChannel.ID

    for {
        messages, err := s.ChannelMessages(dataChannelID, 100, beforeMessageID, "", "")
        if err != nil {
            break
        }

        if len(messages) < 1 {
            break
        }

        beforeMessageID = messages[len(messages) - 1].ID

        for _, dataMessage := range messages {
            userData, err = UserDataFromDiscordDataStoreMessage(dataMessage)
            if err != nil {
                s.ChannelMessageSend(m.ChannelID, err.Error())
                return nil, "", err
            }

            if userData.UserID == userID {
                return userData, dataMessage.ID, nil
            }
        }
    }

    return nil, "n/a", errors.New("No user data found")
}

func UserDataFromDiscordDataStoreMessage(m *dgo.Message) (*UserData, error) {
    messageData := m.ContentWithMentionsReplaced()

    var userData UserData
    err := json.Unmarshal([]byte(messageData), &userData)
    if err != nil {
        return nil, err
    }

    return &userData, nil
}

func CheckCommand(s *dgo.Session, m *dgo.Message) (string, bool) {
	if CheckForSelf(s, m) {
		return "", true
	}

	if !strings.HasPrefix(m.Content, DiscordBotPrefix) {
		return "", true
	}

	newString := strings.Replace(m.Content, DiscordBotPrefix, "", 1)

	return newString, false
}

func CheckForSelf(s *dgo.Session, m *dgo.Message) bool {
    if m.Author.ID == s.State.User.ID {
        return true
    }
    return false
}

func GetMessageGuild(s *dgo.Session, m *dgo.Message) (*dgo.Guild, error) {
	for _, guild := range s.State.Guilds {
		if strings.EqualFold(guild.ID, m.GuildID) {
			return guild, nil
		}
	}
	return nil, errors.New("Unable to find guild for message")
}

func GetChannelFromGuild(channelName string, guild *dgo.Guild) (*dgo.Channel, error) {
	for _, channel := range guild.Channels {
		if strings.EqualFold(channel.Name, channelName) {
			return channel, nil
		}
	}
	return nil, errors.New("Unable to find channel")
}

func CheckForRole(roleName string, s *dgo.Session, m *dgo.Message) bool {
	messageGuild, err := GetMessageGuild(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Couldn't find guild with error: %s", err))
	}

	roleID := ""
	for _, role := range messageGuild.Roles {
		if strings.EqualFold(role.Name, roleName) {
			roleID = role.ID
		}
	}

	if len(roleID) < 1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Couldn't find role: %s", roleName))
		return false
	}

	for _, memberRoleID := range m.Member.Roles {
		if strings.EqualFold(memberRoleID, roleID) {
			return true
		}
	}

	return false
}
