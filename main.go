package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	if len(Token) == 0 {
		Token = os.Getenv("TOKEN")
	}
}

func main() {
	fmt.Println(Token)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	// dg.AddHandler(messageCreate)
	dg.AddHandler(voiceStateUpdate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

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
	fmt.Println("Update")
	vc, err := session.Channel(event.ChannelID)

	if err != nil {
		fmt.Println("Channel not found :", err)
		return;
	}
	
	g, err := session.Guild(event.GuildID)

	if err != nil {
		fmt.Println("Guild not found :", err)
		return;
	}

	if vc.ParentID == "775393101975388232" && vc.ID != "775399800207835140" {
			// g.VoiceStates
			if getUserAmountByChannelId(g.VoiceStates, vc.ID) == 0 {

			}		
	} else if vc.ID == "775399800207835140" {
		var newChannel discordgo.GuildChannelCreateData;
		newChannel.Name = "Voice Channel"
		newChannel.ParentID = vc.ParentID
		newChannel.Type = 2 // Voice Channel

		nc, err := session.GuildChannelCreateComplex(vc.GuildID, newChannel)

		if err != nil {
			fmt.Println("Channel could not be created:", err)
			return;
		}

		session.GuildMemberMove(vc.GuildID, event.UserID, &nc.ID)
	}

	// fmt.Println(s.StateEnabled)
}

func getUserAmountByChannelId(states []*discordgo.VoiceState, cid string) (int) {
	var amount int = 0
	for _, state := range states {
		if state.ChannelID == cid {
			amount++
		}
	}
	return amount
}