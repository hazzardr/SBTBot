package discord

import (
	"context"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
	"github.com/disgoorg/snowflake/v2"
)

type SlashCommand struct {
	Metadata   discord.ApplicationCommandCreate
	Path       string
	HandleFunc handler.SlashCommandHandler
}

var staticCommands = []SlashCommand{
	pingCommand,
}

type Bot struct {
	client *bot.Client
}

func NewBot(discordToken string) (*Bot, error) {
	b := &Bot{}
	r := initializeEventListeners(b)

	client, err := disgo.New(discordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
			),
			gateway.WithPresenceOpts(
				gateway.WithListeningActivity("!"),
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
			),
		),
		bot.WithEventListeners(r),
	)
	if err != nil {
		return nil, err
	}
	b.client = client
	return b, nil
}

func initializeEventListeners(b *Bot) *handler.Mux {
	r := handler.New()
	r.Use(middleware.Logger)
	for _, cmd := range staticCommands {
		r.SlashCommand(cmd.Path, cmd.HandleFunc)
	}

	submitModalCMD := getLaunchSubmitModalCommand(b)
	r.SlashCommand(submitModalCMD.Path, submitModalCMD.HandleFunc)

	return r
}

// Start starts the bot, initializing the discord client against the gateway.
func (b *Bot) Start() error {
	if err := b.client.OpenGateway(context.Background()); err != nil {
		return err
	}
	return nil
}

func (b *Bot) GracefulShutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	b.client.Close(ctx)
}

func (b *Bot) SyncCommands() error {
	creates := make([]discord.ApplicationCommandCreate, 0)
	for _, cmd := range staticCommands {
		creates = append(creates, cmd.Metadata)
	}
	creates = append(creates, getLaunchSubmitModalCommand(b).Metadata)
	err := handler.SyncCommands(b.client, creates, make([]snowflake.ID, 0))
	return err
}
