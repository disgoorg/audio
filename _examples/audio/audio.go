package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/disgoorg/audio"
	"github.com/disgoorg/audio/mp3"
	"github.com/disgoorg/audio/pcm"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

var player audio.Player

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Info("starting up")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentMessageContent|gateway.IntentGuildMessages|gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
		bot.WithEventListenerFunc(func(e *events.GuildMessageCreate) {
			args := strings.Split(e.Message.Content, " ")
			switch args[0] {
			case "pause":
				player.SetPaused(true)
			case "resume":
				player.SetPaused(false)
			case "volume":
				volume, _ := strconv.ParseFloat(args[1], 64)
				player.SetVolume(float32(volume))
			}
		}),
	)
	if err != nil {
		log.Fatal("error creating client: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func play(client bot.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := client.OpenVoice(ctx, guildID, channelID, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err = conn.WaitUntilConnected(ctx); err != nil {
		panic("error waiting for connection: " + err.Error())
	}

	rs, err := http.Get("https://p.scdn.co/mp3-preview/029f4fba66c0b2cfddfe53fc14b95fa2982e423a")
	if err != nil {
		panic("error getting audio: " + err.Error())
	}

	mp3Provider, writer, err := mp3.NewPCMFrameProvider(nil)
	if err != nil {
		panic("error creating mp3 provider: " + err.Error())
	}

	opusProvider, err := pcm.NewOpusProvider(nil, mp3Provider)
	if err != nil {
		panic("error creating opus provider: " + err.Error())
	}

	conn.SetOpusFrameProvider(opusProvider)

	defer rs.Body.Close()
	if _, err = io.Copy(writer, rs.Body); err != nil {
		panic("error reading audio: " + err.Error())
	}
}
