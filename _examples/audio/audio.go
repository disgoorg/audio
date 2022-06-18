package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgoplayer"
	"github.com/disgoorg/disgoplayer/audio/mp3"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Info("starting up")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
	)
	if err != nil {
		log.Fatal("error creating client: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func play(client bot.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connection, err := client.ConnectVoice(ctx, guildID, channelID, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err = connection.WaitUntilConnected(ctx); err != nil {
		panic("error waiting for connection: " + err.Error())
	}

	rs, err := http.Get("https://p.scdn.co/mp3-preview/ee121ca281c629bb4e99c33d877fe98fbb752289?cid=774b29d4f13844c495f206cafdad9c86")
	if err != nil {
		panic("error getting audio: " + err.Error())
	}

	provider, writer, err := mp3.NewMP3PCMFrameProvider(nil)
	if err != nil {
		panic("error creating audio provider: " + err.Error())
	}

	opusProvider, err := disgoplayer.NewPCMOpusProvider(nil, provider)
	if err != nil {
		panic("error creating opus provider: " + err.Error())
	}

	connection.SetOpusFrameProvider(opusProvider)

	println("voice: ready")

	defer rs.Body.Close()
	if _, err = io.Copy(writer, rs.Body); err != nil {
		panic("error reading audio: " + err.Error())
	}
}
