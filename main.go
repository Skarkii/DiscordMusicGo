package main

import (
	"DiscordMusicGo/discordapi"
	"encoding/json"
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
			var payload discordapi.GatewayPayload
			if err := s.GetPayload(&payload); err != nil {
				fmt.Println("Error reading from Discord")
			}

			//fmt.Printf("Received payload: %v\n", payload)

			if payload.Type == "MESSAGE_CREATE" {
				var msg discordapi.MessageCreate
				data, _ := json.Marshal(payload.Data)
				if err := json.Unmarshal(data, &msg); err != nil {
					fmt.Printf("Error unmarshaling message: %v\n", err)
					continue
				}
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
			}

		}
	}()

	if s.Bot.Name != "" {
		fmt.Printf("Bot \"%s\", is now running! Press Ctrl+C to exit.\n", s.Bot.Name)
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}

	err = s.Exit()
	if err != nil {
		log.Fatal("Failed to gracefully shutdown:", err)
	}
	fmt.Println("Bot shutdown gracefully")
}
