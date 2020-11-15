package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Discord Token varibale
var Token string

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	if len(Token) == 0 {
		Token = os.Getenv("TOKEN")
	}
}

func main() {
	loadStorage()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Add the voice state update handler
	dg.AddHandler(voiceStateUpdate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// func channeUpdate(session *discordgo.Session, event *discordgo.ChannelUpdate) {
// 	if event.Type != 2 {
// 		return
// 	}

// 	gs, gsExist := storage.Guilds[event.GuildID]

// 	if !gsExist {
// 		return
// 	}

// 	if event.ParentID == gs.ChannelCategory {
// 		// TODO get creator
// 		// storage.ChannelNames[CREAROT ID] = event.Name
// 		// storage.save
// 	}
// }

func voiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	storedGuild, storedGuildExist := storage.Guilds[event.GuildID]
	if !storedGuildExist {
		return
	}

	g, err := session.State.Guild(event.GuildID)

	if err != nil {
		fmt.Println("Guild not found :", err)
		return
	}

	if len(event.ChannelID) == 0 {
		c, err := session.GuildChannels(event.GuildID)

		if err != nil {
			fmt.Println("Could not retrive channel list for guild:", err)
		}

		for _, channel := range c {
			// Check only for voice channels, and do not delete our creation channel
			if channel.Type == 2 && channel.ParentID == storedGuild.ChannelCategory && channel.ID != storedGuild.CreationChannel {
				if getUserAmountByChannelId(g.VoiceStates, channel.ID) == 0 {
					fmt.Println("Deleted unused channel: ", channel.ID)
					session.ChannelDelete(channel.ID)
				}
			}
		}
		return
	}

	vc, err := session.Channel(event.ChannelID)

	if err != nil {
		fmt.Println("Channel not found :", err)
		return
	}

	if vc.ID == storedGuild.CreationChannel {
		var newChannel discordgo.GuildChannelCreateData
		if len(storage.ChannelNames[event.UserID]) > 0 {
			newChannel.Name = storage.ChannelNames[event.UserID]
		} else {
			newChannel.Name = "Voice Channel"
		}

		newChannel.ParentID = vc.ParentID
		newChannel.Type = 2 // Voice Channel

		nc, err := session.GuildChannelCreateComplex(vc.GuildID, newChannel)

		if err != nil {
			fmt.Println("Channel could not be created:", err)
			return
		}

		session.GuildMemberMove(vc.GuildID, event.UserID, &nc.ID)
	}
}

func getUserAmountByChannelId(states []*discordgo.VoiceState, cid string) int {
	var amount int = 0
	for _, state := range states {
		if state.ChannelID == cid {
			amount++
		}
	}
	return amount
}
