package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"strings"

	"github.com/antonio-petrillo/dixieflatline/history"

	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var server = flag.String("s", "irc.libera.chat", "address of the irc server to connect to")
var port = flag.String("p", "6697", "port to connect to")
var password = flag.String("pass", "s3c43t", "password to use during connection")

var nick = flag.String("n", "dixieflatline", "nickname for the bot")
var user = flag.String("u", "McCoy Pauley", "username for the bot")

type joinFlag []string

func (jf *joinFlag) String() string {
	return strings.Join(*jf, ", ")
}

func (jf *joinFlag) Set(value string) error {
	for _, channel := range strings.Split(value, ",") {
		name := strings.TrimSpace(channel)
		if !strings.HasPrefix(name, "#") && len(name) > 0 {
			name = "#" + name
		}
		*jf = append(*jf, name)
	}
	return nil
}

func main() {
	fmt.Printf("Bot starting\n")

	var joinChannels joinFlag
	flag.Var(&joinChannels, "j", "channels to join automatically")

	flag.Parse()
	channels := func(bot *hbot.Bot) {
		// bot.Channels = []string{"#test"}
		bot.Channels = []string(*&joinChannels)
	}
	user := func(bot *hbot.Bot) {
		bot.Realname = *user
	}
	// better setup this options
	tlsOption := func(bot *hbot.Bot) {
		bot.SSL = true
		bot.SASL = false
		bot.TLSConfig = tls.Config{
			InsecureSkipVerify: true,
		}
		bot.Password = *password
	}

	irc, err := hbot.NewBot(*server + ":" + *port, *nick, tlsOption, channels, user)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(history.HistorySaveTrigger())
	irc.AddTrigger(history.HistoryShowTrigger())
	irc.Logger.SetHandler(log.StdoutHandler)

	irc.Run()
	fmt.Println("Bot shutting down.")
}
