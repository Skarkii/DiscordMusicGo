package main

import (
	"DiscordMusicGo/discordapi"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_TOKEN")
	intents := discordapi.IntentGuildMessages | discordapi.IntentDirectMessages | discordapi.IntentGuildVoiceStates | discordapi.IntentGuilds | discordapi.IntentMessageContent

	s, err := discordapi.New(Token, intents)
	if err != nil {
		log.Fatal("Error creating Discord session")
	}

	go func() {
		for {
			msgtype, msg, err := s.GetMessage()
			if err != nil {
				fmt.Println("Error reading from Discord", err)
				continue
			}

			//fmt.Printf("msgtype: %s, msg: %s\n", msgtype, msg)

			if msgtype == "MESSAGE_CREATE" {
				if msg.Author.ID == s.Bot.ID {
					continue
				}

				fmt.Printf("AUTHOR: %s MSG: %s\n", msg.Author.Name, msg.Content)

				if msg.Content == "-ping" {
					err := s.SendMessage(msg.ChannelID, "Pong!")
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err)
					}
				}
				if msg.Content == "-play" {

					s.ConnectToVoice(msg.GuildID, msg.Author.ID)
				}
			}

		}
	}()

	fmt.Printf("Bot \"%s\", is now running! Press Ctrl+C to exit.\n", s.Bot.Name)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = s.Exit()
	if err != nil {
		log.Fatal("Failed to gracefully shutdown:", err)
	}
	fmt.Println("Bot shutdown gracefully")
}
