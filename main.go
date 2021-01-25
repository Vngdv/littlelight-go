package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// token is used to store the discord token
var token string

func init() {

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if len(token) == 0 {
		token = os.Getenv("TOKEN")
	}
}

func main() {
	// loadStorage()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
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

// Caching stuff would be more efficient but i expect this bot to not be online 24/7
func voiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	// Get the Guild or return. It is only possible in Guilds to create Voice Channels
	g, err := session.State.Guild(event.GuildID)
	if err != nil {
		fmt.Println("Guild not found :", err)
		return
	}
	// Delete channels or rename the last one

	// Retrive the channes of this guild
	c, err := session.GuildChannels(event.GuildID)
	if err != nil {
		fmt.Println("Could not retrive channel list for guild:", err)
	}

	categorys := []*discordgo.Channel{}
	emptyChannels := []*discordgo.Channel{}

GuildChannelLookup:
	for _, channel := range c {
		if channel.Type == discordgo.ChannelTypeGuildCategory && strings.Contains(channel.Name, "🏴") {
			categorys = append(categorys, channel)
			continue
		}

		// Check if it is really a voice channel or continue if it is not
		if channel.Type != discordgo.ChannelTypeGuildVoice {
			continue
		}

		// Do stuff only if the category is okay we skip the current iteration here if the channel has not the right parent
		for _, cat := range categorys {
			if cat.ID != channel.ParentID {
				continue GuildChannelLookup
			}
		}

		// Add empty channels to list
		if UserCount(g.VoiceStates, channel.ID) == 0 {
			emptyChannels = append(emptyChannels, channel)
		}
	}

	// Either create or delete channels
	if len(emptyChannels) > 1 {
		for i, emptyChannel := range emptyChannels {
			if i == len(emptyChannels)-1 {
				fmt.Println("Did not delete channel: ", emptyChannel.Name)
				continue
			}
			fmt.Println("Deleted ", emptyChannel.Name)
			session.ChannelDelete(emptyChannel.ID)
		}
	} else if len(categorys) > 0 {

		var newChannel discordgo.GuildChannelCreateData

		switch rand.Intn(5) {
		case 0:
			newChannel.Name = "Voice Channel"
		case 1:
			newChannel.Name = "Edge Channel"
		case 2:
			newChannel.Name = "🎈 Party Room"
		case 3:
			newChannel.Name = "Kuschel Ecke"
		case 4:
			newChannel.Name = "Another one"
		case 5:
			newChannel.Name = "Just do it"
		}

		newChannel.ParentID = categorys[0].ID
		newChannel.Type = discordgo.ChannelTypeGuildVoice

		_, err := session.GuildChannelCreateComplex(g.ID, newChannel)
		if err != nil {
			fmt.Println("Channel creation failed :", err)
			return
		}

	}
}

// UserCount returns the amount of users that are currently in a channel.
func UserCount(states []*discordgo.VoiceState, cid string) int {
	var amount int = 0
	for _, state := range states {
		if state.ChannelID == cid {
			amount++
		}
	}
	return amount
}
