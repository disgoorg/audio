package main

import (
	"context"
	"fmt"
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
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

var player *TrackPlayer

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Info("starting up")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentMessageContent|gateway.IntentGuildMessages|gateway.IntentGuildVoiceStates)),
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

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func start(client bot.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := client.OpenVoice(ctx, guildID, channelID, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err = conn.WaitUntilConnected(ctx); err != nil {
		panic("error waiting for connection: " + err.Error())
	}

	player = &TrackPlayer{
		queue: []string{
			"https://p.scdn.co/mp3-preview/029f4fba66c0b2cfddfe53fc14b95fa2982e423a",
			"https://p.scdn.co/mp3-preview/53d1fc1d65f13679db03cf7ecb7500212238d998",
			"https://p.scdn.co/mp3-preview/b34cc4a94716e02111c1247fbf963de4ff7b859f",
		},
		conn:   conn,
		client: client,
	}

	player.Player, err = audio.NewPlayer(func() pcm.FrameProvider {
		return player.provider
	}, player)
	if err != nil {
		panic("error creating player: " + err.Error())
	}

	conn.SetOpusFrameProvider(player)
}

type TrackPlayer struct {
	audio.Player
	queue    []string
	provider pcm.FrameProvider
	conn     voice.Conn
	client   bot.Client
}

func (p *TrackPlayer) next() {
	if len(p.queue) == 0 {
		_ = p.client.CloseVoice(context.Background(), p.conn.GuildID())
		return
	}
	var track string
	track, p.queue = p.queue[0], p.queue[1:]

	time.Sleep(time.Second * 2)

	rs, err := http.Get(track)
	if err != nil {
		return
	}
	defer rs.Body.Close()

	var w io.Writer
	p.provider, w, err = mp3.NewPCMFrameProvider(nil)
	if err != nil {
		panic("error creating mp3 provider: " + err.Error())
		return
	}
	_, _ = io.Copy(w, rs.Body)
}

func (p *TrackPlayer) OnPause(player audio.Player) {
	println("paused")
}

func (p *TrackPlayer) OnResume(player audio.Player) {
	println("resume")
}

func (p *TrackPlayer) OnStart(player audio.Player) {
	println("start")
}

func (p *TrackPlayer) OnEnd(player audio.Player) {
	println("end")
	p.next()
}

func (p *TrackPlayer) OnError(player audio.Player, err error) {
	fmt.Println("error: ", err)
}

func (p *TrackPlayer) OnClose(player audio.Player) {
	println("close")
}
