package history

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/antonio-petrillo/dixieflatline/db"
	"github.com/antonio-petrillo/dixieflatline/message"
	"github.com/antonio-petrillo/dixieflatline/utils"

	hbot "github.com/whyrusleeping/hellabot"
)

const (
	MAX_HIST = 25
)

func HistorySaveTrigger(conn *sql.DB) *hbot.Trigger {
	return &hbot.Trigger{
		Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
			return m.Command == "PRIVMSG" && !(strings.HasPrefix(m.Trailing(), "!ignore") || strings.HasPrefix(m.Trailing(), "!history") || strings.HasPrefix(m.Trailing(), "!i") || strings.HasPrefix(m.Trailing(), "!h"))
		},
		Action: func(bot *hbot.Bot, m *hbot.Message) bool {
			size := len(m.Params)
			if size > 1 { // if size <= 1 then it is only trailing and I have no channel to store
				channels := utils.GetChannelsFromParams(m.Params)
				for _, channel := range channels {
					db.StoreHistoryEntry(conn, channel, message.HistoryEntry{From: m.From, Msg: m.Trailing()})
				}
			}
			return false
		},
	}
}

func HistoryShowTrigger(conn *sql.DB) *hbot.Trigger {
	return &hbot.Trigger{
		Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
			return m.Command == "PRIVMSG" && (strings.HasPrefix(m.Trailing(), "!history") || strings.HasPrefix(m.Trailing(), "!h "))
		},
		Action: func(bot *hbot.Bot, m *hbot.Message) bool {
			size := len(m.Params)
			if size > 1 {
				channels := utils.GetChannelsFromParams(m.Params)
				for _, channel := range channels {
					messages, err := db.RetrieveHistoryEntries(conn, channel, MAX_HIST)
					if err == nil {
						for i := len(messages) - 1; i >= 0; i-- {
							bot.Send(fmt.Sprintf("PRIVMSG %s [in %s by %s]> %s", m.From , channel, messages[i].From, messages[i].Msg))
						}
					} else {
						bot.Send(fmt.Sprintf("PRIVMSG %s :Retrieve history failed for %s", m.From, channel))
					}
				}
			}
			return false
		},
	}
}
