package history

import (
	"fmt"
	"github.com/antonio-petrillo/dixieflatline/utils"
	hbot "github.com/whyrusleeping/hellabot"
	"strings"
)

const (
	MAX_HIST = 100
)

type HistoryEntry struct {
	From string
	Msg  string
}

// in memory database (pretty bad), also no locks guards the maps!!!
var channelsHistories map[string][]HistoryEntry = map[string][]HistoryEntry{}

// unexported var
var historySaveCommand = &hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && !(strings.HasPrefix(m.Trailing(), "!ignore") || strings.HasPrefix(m.Trailing(), "!history"))
	},
	Action: func(bot *hbot.Bot, m *hbot.Message) bool {
		size := len(m.Params)
		if size > 1 { // if size <= 1 then it is only trailing and I have no channel to store
			channels := utils.GetChannelsFromParams(m.Params)
			for _, channel := range channels {
				channelHistory := channelsHistories[channel]
				if len(channelHistory) > MAX_HIST {
					channelHistory = channelHistory[1:]
				}
				channelHistory = append(channelHistory, HistoryEntry{From: m.From, Msg: m.Trailing()})
				channelsHistories[channel] = channelHistory
			}
		}

		return false
	},
}

// unexported var
var historyShowCommand = &hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Trailing(), "!history")
	},
	Action: func(bot *hbot.Bot, m *hbot.Message) bool {
		size := len(m.Params)
		if size > 1 {
			channels := utils.GetChannelsFromParams(m.Params)
			for _, channel := range channels {
				history := channelsHistories[channel]
				if len(history) > 0 {
					for _, message := range history {
						bot.Send(fmt.Sprintf("PRIVMSG %s [in %s by %s]> %s", m.From , channel, message.From, message.Msg))
					}
				} else {
				 	bot.Send(fmt.Sprintf("PRIVMSG %s :no prev history", m.From))
				}
			}
		}

		return false
	},
}

func HistorySaveTrigger() *hbot.Trigger {
	return historySaveCommand
}

func HistoryShowTrigger() *hbot.Trigger {
	return historyShowCommand
}
