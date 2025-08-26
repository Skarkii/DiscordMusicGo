# Discord Music Go

A Discord music bot written in Go with custom Discord API implementation.

## Prerequisites

Before running this bot, make sure you have the following installed:

### Required Software
- **Go 1.25+** - [Download Go](https://golang.org/dl/)

### Discord Bot Setup
1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application and bot
3. Copy the bot token for later use
4. Enable the following bot permissions:
   - Send Messages
   - Connect to Voice Channels
   - Speak in Voice Channels
   - Read Message History

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Skarkii/DiscordMusicGo
   cd DiscordMusicGo
   ```

2. **Install Go dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and add your Discord bot token:
   ```env
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

4. **Run the bot:**
   ```bash
   go run .
   ```

## Building

To build the bot as an executable:

```bash
go build -o discord-music-bot
```