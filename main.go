// var (
// token string = os.Getenv("DISCORD_TOKEN")
// channelIDstring string = os.Getenv("DISCORD_CHANNEL_ID")
// guildIDstring   string = os.Getenv("DISCORD_GUILD_ID")
// channelIDraw, _              = strconv.ParseUint(os.Getenv("DISCORD_CHANNEL_ID"), 10, 64)
// guildIDraw, _                = strconv.ParseUint(os.Getenv("DISCORD_GUILD_ID"), 10, 64)
// channelID       snowflake.ID = snowflake.ID(channelIDraw)
// guildID         snowflake.ID = snowflake.ID(guildIDraw)
// )
//var guilds map[string]*Guild

package main

import (
	"context"
	"encoding/binary"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"

	"github.com/joho/godotenv"
)

// var token string = os.Getenv("DISCORD_TOKEN")
// var channelIDraw, _ = strconv.ParseUint(os.Getenv("DISCORD_CHANNEL_ID"), 10, 64)
// var guildIDraw, _ = strconv.ParseUint(os.Getenv("DISCORD_GUILD_ID"), 10, 64)
// var channelID snowflake.ID = snowflake.ID(channelIDraw)
// var guildID snowflake.ID = snowflake.ID(guildIDraw)
// print("token: ", token)
// print("\nguildID: ", guildID)
// print("\nchannelID: ", channelID)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var token string = os.Getenv("DISCORD_TOKEN")
	print("token: ", token)
	slog.Info("starting up")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	s := make(chan os.Signal, 1)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client(), s)
		}),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		client.Close(ctx)
	}()

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error connecting to gateway", slog.Any("error", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func play(client bot.Client, closeChan chan os.Signal) {
	var channelIDraw, _ = strconv.ParseUint(os.Getenv("DISCORD_CHANNEL_ID"), 10, 64)
	var guildIDraw, _ = strconv.ParseUint(os.Getenv("DISCORD_GUILD_ID"), 10, 64)
	var channelID snowflake.ID = snowflake.ID(channelIDraw)
	var guildID snowflake.ID = snowflake.ID(guildIDraw)

	var channelIDraw2, _ = strconv.ParseUint(os.Getenv("DISCORD_CHANNEL_ID2"), 10, 64)
	var guildIDraw2, _ = strconv.ParseUint(os.Getenv("DISCORD_GUILD_ID2"), 10, 64)
	var channelID2 snowflake.ID = snowflake.ID(channelIDraw)
	var guildID2 snowflake.ID = snowflake.ID(guildIDraw)

	conn := client.VoiceManager().CreateConn(guildID)
	conn := client.VoiceManager().CreateConn(guildID2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := conn.Open(ctx, channelID, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer closeCancel()
		conn.Close(closeCtx)
	}()

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	for i := 1; i < 20; i++ {
		writeOpus(conn.UDP())
	}
	closeChan <- syscall.SIGTERM
}

func writeOpus(w io.Writer) {
	file, err := os.Open("nico.dca")
	if err != nil {
		panic("error opening file: " + err.Error())
	}
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenBuf [4]byte
	for range ticker.C {
		_, err = io.ReadFull(file, lenBuf[:])
		if err != nil {
			if err == io.EOF {
				_ = file.Close()
				return
			}
			panic("error reading file: " + err.Error())
			return
		}

		// Read the integer
		frameLen := int64(binary.LittleEndian.Uint32(lenBuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, file, frameLen)
		if err != nil && err != io.EOF {
			_ = file.Close()
			return
		}
	}
}
