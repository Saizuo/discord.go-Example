package discord

// import (
// 	"strings"
// 	"time"

//     // "fpbot/pkg/utils"
//     // "fpbot/pkg/discord/discord_cmd"
//     fpmodel "fpbot/pkg/model"

// 	dgo "github.com/bwmarrin/discordgo"
//     // db "github.com/replit/database-go"
// )

// func NewBotData() fpmodel.BotData {
//     return fpmodel.BotData {
//         StartTime:    time.Now(),
//         LastRateLimitedCommandTime: time.Now(),
//     }
// }

// func (bd *fpmodel.BotData) handleRegularText(s *dgo.Session, m *dgo.MessageCreate) {
// 	message, failed := utils.CheckCommand(s, m.Message)
// 	if failed || utils.CheckForSelf(s, m.Message) {
//         return
// 	}

//     splitMessage := strings.Split(message, " ")

//     cmd := discord_cmd.NewCommand(s, m.Message, bd, splitMessage)

//     err := cmd.Execute()
//     if err != nil {
//         s.ChannelMessageSend(m.ChannelID, err.Error())
//     }
// }
