package main

import (
	// "main/dgvoice"
	// "main/ytdlp"

	// "fmt"
	"log"
	"os"
	"os/signal"
	// "strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var guilds map[string]*Guild

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_TOKEN")
	//Prefix := os.Getenv("COMMAND_PREFIX")

	session, _ := discordgo.New("Bot " + Token)

	guilds = make(map[string]*Guild)

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

	// If guild exists, pass message. If not, create a new one and handle message
	if guild, exists := guilds[m.GuildID]; exists {
		guild.HandleMessage(s, m)
	} else {
		newGuild := NewGuild(m.GuildID, "-p")
		guilds[m.GuildID] = newGuild
		newGuild.HandleMessage(s, m)
	}
}
