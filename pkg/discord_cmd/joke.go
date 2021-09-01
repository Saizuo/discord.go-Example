package discord_cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "time"

    // fpmodel "fpbot/pkg/model"

    "github.com/spf13/cobra"
    dgo "github.com/bwmarrin/discordgo"
)

const dadJokeBaseURL = "https://icanhazdadjoke.com/"

type jokeCommand struct {
    DiscordCommand
}

func (c *jokeCommand) run() {
    if time.Since(c.BotData.LastRateLimitedCommandTime).Seconds() < 5 {
        c.Session.ChannelMessageSend(c.Message.ChannelID, "Slow down on the jokes!")
        return
    }
    client := &http.Client{}

    req, err := http.NewRequest("GET", dadJokeBaseURL, nil)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Cannot create GET request: %s", err.Error()))
        return
    }
    req.Header.Set("User-Agent", "Friendly Potato Discord Bot (https://github.com/you-win/fpbot)")
    req.Header.Set("Accept", "application/json")

    res, err := client.Do(req)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Error when sending GET request: %s", err.Error()))
        return
    }

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)

    var jsonBody map[string]interface{}
    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Error when unmarshalling response: %s", err.Error()))
        return
    }

    if jsonBody["status"].(float64) != 200 {
        c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Non-200 response code: %s", jsonBody))
        return
    }

    c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("```%s```", jsonBody["joke"].(string)))

    c.BotData.LastRateLimitedCommandTime = time.Now()
}

func NewJokeCommand(s *dgo.Session, m *dgo.Message, b *BotData) *cobra.Command {
    dc := &jokeCommand{
        DiscordCommand: DiscordCommand{
            Session: s,
            Message: m,
            BotData: b,
        },
    }

    c := &cobra.Command{
        Use: "joke",
        Short: "Have the bot tell a joke",
        Run: func(cmd *cobra.Command, args []string){
            dc.run()
        },
    }
    
    modifyUsageFunc(c, s, m)

    return c
}
