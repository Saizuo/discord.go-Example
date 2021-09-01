package chat_games

import (
	"log"
	"strconv"
	"time"
    "strings"

	database "github.com/replit/database-go"
)

const (
    CountUpName = "CountUp"
    CountUpStateName = "CountUpState"
)

type CountUp struct {
	ChatGame
	score      int
	lastNumber int // Be careful of negative numbers
}

func NewCountUp(channelID string) CountUp {
	return CountUp{
		ChatGame: ChatGame{
            CleanupTime: time.Now().Add(time.Hour),
            ChannelID: channelID,
        },
		score:            0,
		lastNumber:       0,
	}
}

func (c *CountUp) GetCleanupTime() time.Time {
    return c.CleanupTime
}

func (c *CountUp) GetChannelID() string {
    return c.ChannelID
}

func (c *CountUp) Play(value string) bool {
	num, err := strconv.Atoi(value)
	if err != nil {
		// Passed a bad number, automatically lose the game
		return false
	}

	if num-c.lastNumber != 1 {
		return false
	}

    c.lastNumber = num
    c.score += 1
    c.CleanupTime = time.Now().Add(time.Hour)
	return true
}

func (c *CountUp) ReadScore() (string, error) {
	dbValue, err := database.Get(CountUpName)
	if err != nil {
		return "0", err
	}

	return dbValue, nil
}

func (c *CountUp) WriteScore() error {
	lastHighScore, err := c.ReadScore()
	if err != nil {
        // We don't return here since this usually means there is no score data
		log.Println(err)
        lastHighScore = "0"
	}
	uintHighScore, err := strconv.Atoi(lastHighScore)
	if err != nil {
		return err
	}

	if c.score > uintHighScore {
		err = database.Set(CountUpName, strconv.Itoa(c.score))
        if err != nil {
            return err
        }
	}

	return nil
}

func (c *CountUp) ReadState() error {
    dbValue, err := database.Get(CountUpStateName)
    if err != nil {
        return err
    }

    splitData := strings.Split(dbValue, ",")
    parsedCleanupTime, err := time.Parse(time.UnixDate, splitData[0])
    if err != nil {
        return err
    }
    c.CleanupTime = parsedCleanupTime

    parsedScore, err := strconv.Atoi(splitData[1])
    if err != nil {
        return err
    }
    c.score = parsedScore

    parsedLastNumber, err := strconv.Atoi(splitData[2])
    if err != nil {
        return err
    }
    c.lastNumber = parsedLastNumber

    return nil
}

func (c *CountUp) SaveState() error {
    var sb strings.Builder
    sb.WriteString(c.CleanupTime.Format(time.UnixDate))
    sb.WriteString(",")
    sb.WriteString(strconv.Itoa(c.score))
    sb.WriteString(",")
    sb.WriteString(strconv.Itoa(c.lastNumber))

    database.Set(CountUpStateName, sb.String())

	return nil
}

func (c *CountUp) Score() interface{} {
    return c.score
}
