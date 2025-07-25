package main

import (
	"fmt"
	"strings"
	"log"
	"github.com/bwmarrin/discordgo"
	"main/dgvoice"
	"main/ytdlp"
)

type Guild struct {
	ID string
	prefix string
	voiceConnection *discordgo.VoiceConnection
}

func NewGuild(ID string, prefix string) *Guild {
	log.Printf("Created new guild with id %s\n", ID)
	return &Guild{
		ID: ID,
		prefix: prefix,
		voiceConnection: nil,
	}
}

func (g *Guild) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, g.prefix) == false {
		log.Printf("[INFO] [%s] : User message had Invalid Prefix", g.ID)
		return
	}
	log.Printf("[INFO] [%s] : %s", g.ID, m.Message.Content)

	voiceChannelID, err := g.getUserVoiceChannelID(s, m.Author.ID)
	if err != nil {
		log.Printf("[ERROR] Error fetching voice state: %v\n", err)
		return
	}

	if voiceChannelID == "" {
		log.Printf("[INFO] : User is not in a voice channel ignoring\n")
		return
	}

	g.voiceConnection, err = s.ChannelVoiceJoin(g.ID, voiceChannelID, false, true)
	if err != nil {
		log.Printf("[ERROR] Error joining voice channel %s: %v\n", voiceChannelID, err)
	}

	url := ytdlp.GetYTDLPCommand(strings.TrimPrefix(m.Message.Content, g.prefix))
	go dgvoice.PlayAudioFile(g.voiceConnection, url.Stdout, make(chan bool))
}

func (g *Guild) getUserVoiceChannelID(s *discordgo.Session, userID string) (string, error) {
	guild, err := s.State.Guild(g.ID)
	if err != nil {
		guild, err = s.Guild(g.ID)
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
