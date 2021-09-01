package discord

import (
	"fpbot/pkg/utils"

	"fmt"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

const (
	maxStrikes                uint    = 5
	maxSecondsBetweenMessages float64 = 5
    resetSeconds float64 = 30
)

type AntiSpam struct {
	spamData *SpamData
}

func NewAntiSpam() AntiSpam {
	return AntiSpam{
		spamData: &SpamData{},
	}
}

type SpamData struct {
	lastChatterID        string
	lastMessage          string
	timeSinceLastMessage time.Time
	strikes              uint
}

func (b AntiSpam) handleSpam(s *dgo.Session, m *dgo.MessageCreate) {
	// Don't ban ourself
	if utils.CheckForSelf(s, m.Message) {
		return
	}

	// I can spam
	if m.Author.Username == "youwin" && m.Author.Discriminator == "5391" {
	    return
	}

	// New chatter
	if b.spamData.lastChatterID != m.Author.ID {
		b.spamData.lastChatterID = m.Author.ID
		b.spamData.lastMessage = m.Message.Content
		b.spamData.timeSinceLastMessage = time.Now()
		b.spamData.strikes = 0

		return
	}

	duration := time.Since(b.spamData.timeSinceLastMessage)
    if duration.Seconds() > resetSeconds {
        b.spamData.strikes = 0
        return
    }

	if b.spamData.lastMessage == m.Message.Content {
		// Same chatter, same message
		b.spamData.addStrike()
	} else if duration.Seconds() < maxSecondsBetweenMessages {
		// Same chatter, different message
		b.spamData.lastMessage = m.Message.Content
		b.spamData.addStrike()
	} else {
		b.spamData.strikes = 0
	}

	if b.spamData.strikes == maxStrikes-2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are 2 strikes away from being banned for spam: %s", m.Author.Username))
	} else if b.spamData.strikes == maxStrikes-1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are 1 strike away from being banned for spam: %s", m.Author.Username))
	} else if b.spamData.strikes >= maxStrikes {
		// Finally ban someone if they accrue too many strikes
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Banning user %s for spam", m.Author.ID))
		s.GuildBanCreateWithReason(m.GuildID, m.Author.ID, "Spam", 1)
	}
}

func (b *SpamData) addStrike() {
	b.timeSinceLastMessage = time.Now()
	b.strikes += 1
}
