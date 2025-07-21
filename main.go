package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"io/ioutil"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_TOKEN")
	Prefix := os.Getenv("COMMAND_PREFIX")

	fmt.Printf("Token=%s, Prefix=%s", Token, Prefix)

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
	// Ignore messages sent if they're not sent from the designate bot channel
	if m.ChannelID != "890646235478904893" {
		return
	}
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	voiceChannelID, err := getUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		log.Printf("Error fetching voice state: %v", err)
		return
	}

	var dgv *discordgo.VoiceConnection

	if voiceChannelID != "" {
		fmt.Printf("User is in channel %s\n", voiceChannelID)
		dgv, err = s.ChannelVoiceJoin(m.GuildID, voiceChannelID, false, true)
		if err != nil {
			log.Printf("Error joining voice channel %s: %v", voiceChannelID, err)
		}
	} else {
		fmt.Printf("User is NOT in channel %s\n", voiceChannelID)
		var reply string
		reply = fmt.Sprintf("Hey %s you are not in any channel!", m.Author.DisplayName())
		_ = reply
		// s.ChannelMessageSend(m.ChannelID, reply)
	}

	files, _ := ioutil.ReadDir("songs")

	for _, f := range files {
		fmt.Println("PlayAudioFile:", f.Name())
		dgvoice.PlayAudioFile(dgv, fmt.Sprintf("%s/%s", "songs", f.Name()), make(chan bool))
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

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", nil
}
