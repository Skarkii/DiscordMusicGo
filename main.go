package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Printf("")
	var Token string
	var err error
	Token = ""
	session, _ := discordgo.New("Bot " + Token)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})
	session.AddHandler(messageCreate)
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds
	session.StateEnabled = true

	// Start
	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	// Gracefully shutdown with Ctrl + C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	voiceChannelID, err := getUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		log.Printf("Error fetching voice state: %v", err)
		return
	}

	if voiceChannelID != "" {
		fmt.Printf("User is in channel %s\n", voiceChannelID)
		_, err := s.ChannelVoiceJoin(m.GuildID, voiceChannelID, false, true)
		if err != nil {
			log.Printf("Error joining voice channel %s: %v", voiceChannelID, err)
		}
	} else {
		fmt.Printf("User is NOT in channel %s\n", voiceChannelID)
		var reply string
		reply = fmt.Sprintf("Hey %s you are not in any channel!", m.Author.DisplayName())
		s.ChannelMessageSend(m.ChannelID, reply)
	}
}

func getUserVoiceChannel(s *discordgo.Session, guildID, userID string) (string, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		guild, err = s.Guild(guildID)
		if err != nil {
			return "", fmt.Errorf("error fetching guild: %v", err)
		}
	}

	for i, vs := range guild.VoiceStates {
		fmt.Printf("%d", i)
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", nil
}
