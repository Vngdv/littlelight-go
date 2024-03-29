package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// token is used to store the discord token
	token string

	channelNames = []string{
		"Voice Channel",
		"🎈 Party Room",
	}
	joinChannelName string
	categoryIdentifier string
	defaultBitrate int

	allowUserChannelNames bool
)



func init() {
	var names string
	flag.StringVar(&names, "n", "", "Channel Names")
	flag.StringVar(&joinChannelName, "j", "📢 Join to own", "Join Channel Name")
	flag.StringVar(&categoryIdentifier, "c", "🎤", "Category Identifier")
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.IntVar(&defaultBitrate, "b", 64, "Channel Bitrate")
	flag.BoolVar(&allowUserChannelNames, "allow-channelnames", true, "Allow custom user channel names")
	flag.Parse()

	if len(names) > 0 {
		channelNames = strings.Split(names, ";")
	}

	println("Found", len(channelNames), "Channel names")
	println("Join Channel Name:", joinChannelName)
	println("Category Identifier:", categoryIdentifier)
	println("Channel Bitrate:", defaultBitrate*1000)
	println("Custom user channel names:", allowUserChannelNames)

	if len(token) == 0 {
		token = os.Getenv("TOKEN")
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Set random seed
	rand.Seed(time.Now().Unix())

	// Add the voice state update handler and set the intents
	dg.AddHandler(voiceStateUpdate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds | discordgo.IntentsDirectMessages)

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

func voiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	// Get the Guild or return. It is only possible in Guilds to create Voice Channels
	g, err := session.State.Guild(event.GuildID)
	if err != nil {
		fmt.Println("Guild not found :", err)
		return
	}

	// Retrive the channes of this guild
	c, err := session.GuildChannels(event.GuildID)
	if err != nil {
		fmt.Println("Could not retrive channel list for guild:", err)
		return
	}

	categorys := []*discordgo.Channel{}
	emptyChannels := []*discordgo.Channel{}

GuildChannelLookup:
	for _, channel := range c {
		if channel.Type == discordgo.ChannelTypeGuildCategory && strings.Contains(channel.Name, categoryIdentifier) {
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
	if len(emptyChannels) > 0 {
		for i, emptyChannel := range emptyChannels {
			if i == len(emptyChannels)-1 {
				fmt.Println("Did not delete channel: ", emptyChannel.Name)

				// Rename the channel if neede
				if emptyChannel.Name != joinChannelName {
					// Rename channel
					var channelEdit discordgo.ChannelEdit
					channelEdit.Name = joinChannelName

					session.ChannelEditComplex(emptyChannel.ID, &channelEdit)
				}
				continue
			}
			fmt.Println("Deleted ", emptyChannel.Name)
			session.ChannelDelete(emptyChannel.ID)
		}
	} else if len(emptyChannels) == 0 {

		var newChannel discordgo.GuildChannelCreateData

		// Random channel name
		newChannel.Name = joinChannelName

		newChannel.ParentID = categorys[0].ID
		newChannel.Type = discordgo.ChannelTypeGuildVoice

		if defaultBitrate != 64*1000 {
			newChannel.Bitrate = defaultBitrate * 1000
		}

		_, err := session.GuildChannelCreateComplex(g.ID, newChannel)
		if err != nil {
			fmt.Println("Channel creation failed :", err)
			return
		}
	}

	// ---------------- Channel owning system ----------------

	guildMember, _ := session.GuildMember(event.GuildID, event.UserID)

	// Fully ignore bot users
	if guildMember.User.Bot {
		return
	}

	// Check if the user joind in an empty channel that has now one user
	channel, err := session.State.Channel(event.ChannelID)
	if err != nil {
		return
	}
	parent, err := session.State.Channel(channel.ParentID)
	if err != nil {
		return
	}

	if UserCount(g.VoiceStates, event.ChannelID) == 1 && strings.Contains(parent.Name, categoryIdentifier) {
		var channelEdit discordgo.ChannelEdit

		// Only get the custom channel names if we really need them
		if(allowUserChannelNames) {
			// Get last message
			userChannel, _ := session.UserChannelCreate(event.UserID)
			messages, _ := session.ChannelMessages(userChannel.ID, 1, "", "", "")

			// Use custom name for channel if provided
			if len(messages) > 0 {
				channelEdit.Name = messages[0].Content
			} else {
				channelEdit.Name = channelNames[rand.Intn(len(channelNames))]
			}
		} else {
			channelEdit.Name = channelNames[rand.Intn(len(channelNames))]
		}

		_, err := session.ChannelEditComplex(event.ChannelID, &channelEdit)
		if err != nil {
			fmt.Println("Channel edit failed :", err)
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
