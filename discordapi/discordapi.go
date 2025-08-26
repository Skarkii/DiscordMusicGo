/*
Custom Discord API using websockets. Main focus here is parallel voice channel connections
*/
package discordapi

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	Token   string
	conn    *websocket.Conn
	intents int
}

//func httpRequest() {
//	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
//	req.Header.Set("Authorization", "Bot "+token)
//	req.Header.Set("Content-Type", "application/json")
//	if err != nil {
//		log.Fatalf("Could not send request")
//	}
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatalf("Could not send request")
//		return s, nil
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalf("Could not send request")
//	}
//
//	fmt.Printf("Response from Discord: %s\n", string(body))
//}

const (
	IntentGuilds                      = 1 << 0
	IntentGuildMembers                = 1 << 1
	IntentGuildModeration             = 1 << 2
	IntentGuildExpressions            = 1 << 3
	IntentGuildIntegrations           = 1 << 4
	IntentGuildWebhooks               = 1 << 5
	IntentGuildInvites                = 1 << 6
	IntentGuildVoiceStates            = 1 << 7
	IntentGuildPresences              = 1 << 8
	IntentGuildMessages               = 1 << 9
	IntentGuildMessageReactions       = 1 << 10
	IntentGuildMessageTyping          = 1 << 11
	IntentDirectMessages              = 1 << 12
	IntentDirectMessageReactions      = 1 << 13
	IntentDirectMessageTyping         = 1 << 14
	IntentMessageContent              = 1 << 15
	IntentGuildScheduledEvents        = 1 << 16
	IntentAutoModerationConfiguration = 1 << 20
	IntentAutoModerationExecution     = 1 << 21
	IntentGuildMessagePolls           = 1 << 24
	IntentDirectMessagePolls          = 1 << 25
)

const (
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	OpDispatch           = 0    //	Receive	An event was dispatched.
	OpHeartbeat          = 1    //	Send/Receive	Fired periodically by the client to keep the connection alive.
	OpIdentify           = 2    //	Send	Starts a new session during the initial handshake.
	OpPresence           = 3    //	Update	Send	Update the client's presence.
	OpVoice              = 4    //	State Update	Send	Used to join/leave or move between voice channels.
	OpResume             = 6    //	Send	Resume a previous session that was disconnected.
	OpReconnect          = 7    //	Receive	You should attempt to reconnect and resume immediately.
	OpRequest            = 8    //	Guild Members	Send	Request information about offline guild members in a large guild.
	OpInvalid            = 9    //	Session	Receive	The session has been invalidated. You should reconnect and identify/resume accordingly.
	OpHello              = 10   //	Receive	Sent immediately after connecting, contains the heartbeat_interval to use.
	OpHeartbeatAck       = 11   //	ACK	Receive	Sent in response to receiving a heartbeat to acknowledge that it has been received.
	OpRequestSoundboards = 31   //	Soundboard Sounds	Send	Request information about soundboard sounds in a set of guilds.
	OpClose              = 1000 //  Closes the connection
)

const (
	gateway = "wss://gateway.discord.gg/?v=10&encoding=json"
)

type MessageCreate struct {
	Content string `json:"content"`
	Author  struct {
		ID string `json:"id"`
	} `json:"author"`
	ChannelID string `json:"channel_id"`
}

type GatewayPayload struct {
	Op   int         `json:"op"`
	Data interface{} `json:"d"`
	Seq  int         `json:"s,omitempty"`
	Type string      `json:"t,omitempty"`
}

type Identify struct {
	Token      string             `json:"token"`
	Intents    int                `json:"intents"`
	Properties IdentifyProperties `json:"properties"`
}

type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

func (s Session) Exit() error {
	err := s.disconnect()
	if err != nil {
		log.Printf("Could not disconnect!: %v", err)
	}
	err = s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s Session) disconnect() error {
	// Identifies and connects the bot
	disc := GatewayPayload{
		Op:   OpClose,
		Data: nil,
	}
	if err := s.conn.WriteJSON(disc); err != nil {
		return err
	}
	return nil
}

func (s Session) GetPayload(d *GatewayPayload) error {
	return s.conn.ReadJSON(&d)
}

func New(token string, intents int) (Session, error) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(gateway, nil)

	if err != nil {
		log.Fatalf("Failed to establish connection to Discord: %v", err)
	}

	s := Session{
		token,
		conn,
		intents,
	}
	fmt.Println("Intents: ", s.intents)
	// Identifies and connects the bot
	identify := GatewayPayload{
		Op: OpIdentify,
		Data: Identify{
			Token:   token,
			Intents: s.intents,
			Properties: IdentifyProperties{
				OS:      "Windows 11",
				Browser: "DiscordMusicGo",
				Device:  "DiscordMusicGo",
			},
		},
	}
	if err := conn.WriteJSON(identify); err != nil {
		log.Printf("Error sending IDENTIFY: %v\n", err)
	}

	var heartbeatInterval int
	var payload GatewayPayload
	if err := conn.ReadJSON(&payload); err != nil {
		return s, err
	}

	if payload.Op != OpHello {
		return s, errors.New("Invalid starting operator retrieved!")
	}

	fmt.Printf("Retrieved Ack from Identify, starting heartbeat\n")
	data := payload.Data.(map[string]interface{})
	heartbeatInterval = int(data["heartbeat_interval"].(float64))
	go startHeartbeat(conn, heartbeatInterval)

	return s, nil
}

func startHeartbeat(conn *websocket.Conn, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		payload := GatewayPayload{Op: OpHeartbeat, Data: nil}
		if err := conn.WriteJSON(payload); err != nil {
			fmt.Printf("Error sending heartbeat: %v\n", err)
			return
		}
	}
}
