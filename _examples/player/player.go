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

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/disgoplayer"
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
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentMessageContent|discord.GatewayIntentGuildMessages|discord.GatewayIntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go start(e.Client())
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

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

type player struct {
	p        disgoplayer.Player
	queue    []string
	provider disgoplayer.PCMFrameProvider
}

func newPlayer(p disgoplayer.Player) *player {
	return &player{
		p: p,
		queue: []string{
			"https://p.scdn.co/mp3-preview/029f4fba66c0b2cfddfe53fc14b95fa2982e423a",
			"https://p.scdn.co/mp3-preview/53d1fc1d65f13679db03cf7ecb7500212238d998",
			"https://p.scdn.co/mp3-preview/b34cc4a94716e02111c1247fbf963de4ff7b859f",
		},
	}
}

func start(client bot.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := client.ConnectVoice(ctx, guildID, channelID, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err = conn.WaitUntilConnected(ctx); err != nil {
		panic("error waiting for connection: " + err.Error())
	}

	player, err = disgoplayer.NewPlayer(func() disgoplayer.PCMFrameProvider {
		return mp3Provider
	}, func(err error) {
		if err == disgoplayer.OpusProviderClosed {
			conn.Close()
			return
		}
		if err == io.EOF {
			playNext(client, conn)
			return
		}
		return
	})
	if err != nil {
		panic("error creating player: " + err.Error())
	}

	conn.SetOpusFrameProvider(player)
}

func playNext(client bot.Client, conn voice.Connection) {
	if len(queue) == 0 {
		_ = client.DisconnectVoice(context.Background(), conn.GuildID())
		return
	}
	var track string
	track, queue = queue[0], queue[1:]

	time.Sleep(time.Second * 2)

	rs, err := http.Get(track)
	if err != nil {
		return
	}
	defer rs.Body.Close()

	var w io.Writer
	mp3Provider, w, err = disgoplayer.NewMP3PCMFrameProvider(nil)
	if err != nil {
		panic("error creating mp3 provider: " + err.Error())
		return
	}
	_, _ = io.Copy(w, rs.Body)
}
