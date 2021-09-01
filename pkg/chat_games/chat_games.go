package chat_games

import "time"

type ChatGameData struct {
    Data *ChatGame
}

type ChatGame struct {
    CleanupTime time.Time
    ChannelID string
    PlayableGame
}

type PlayableGame interface {
    GetCleanupTime() time.Time
    GetChannelID() string
    Play(string) bool
    ReadScore() (string, error)
    WriteScore() error
    ReadState() error
    SaveState() error
    Score() interface{}
}
