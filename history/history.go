package history

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/antonio-petrillo/dixieflatline/db"
	"github.com/antonio-petrillo/dixieflatline/message"
	"github.com/antonio-petrillo/dixieflatline/utils"

	hbot "github.com/whyrusleeping/hellabot"
)

const (
	MAX_HIST = 25
)

var commands map[string]struct{} = map[string]struct{}{
	"!h":       struct{}{},
	"!history": struct{}{},
	"!i":       struct{}{},
	"!ignore":  struct{}{},
	"!calc":  struct{}{},
}

func HistorySaveTrigger(conn *sql.DB) *hbot.Trigger {
	return &hbot.Trigger{
		Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
			if m.Command != "PRIVMSG" || len(m.Trailing()) == 0 {
				return false
			}

			maybeCommand := strings.Split(strings.TrimSpace(m.Trailing()), " ")[0]
			_, ok := commands[maybeCommand]
			return !ok
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
			if m.Command != "PRIVMSG" || len(m.Trailing()) == 0 {
				return false
			}
			cmd := strings.Split(strings.TrimSpace(m.Trailing()), " ")[0]
			return cmd == "!history" || cmd == "!h"
		},
		Action: func(bot *hbot.Bot, m *hbot.Message) bool {
			size := len(m.Params)
			if size > 1 {
				numMsg := MAX_HIST
				params := utils.SplitTrailingBySpace(m.Trailing())
				for i := 1; i < len(params)-1; i++ {
					if params[i] == "-l"|| params[i] == "-limit" || params[i] == "--limit" {
						conv, err := strconv.Atoi(params[i+1])
						if err == nil {
							numMsg = conv
						}
					}
				}

				channels := utils.GetChannelsFromParams(m.Params)
				for _, channel := range channels {
					messages, err := db.RetrieveHistoryEntries(conn, channel, numMsg)
					if err == nil {
						for i := len(messages) - 1; i >= 0; i-- {
							bot.Send(fmt.Sprintf("PRIVMSG %s [in %s by %s]> %s", m.From, channel, messages[i].From, messages[i].Msg))
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
