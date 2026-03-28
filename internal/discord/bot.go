package discord

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/hazzardr/sbtbot/internal/sbtb"
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
	db     *sbtb.DB
	client *bot.Client
}

func NewBot(discordToken string, dbPath string) (*Bot, error) {
	db, err := sbtb.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	b := &Bot{db: db}
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
	r.Use(loggingMiddleware)
	for _, cmd := range staticCommands {
		r.SlashCommand(cmd.Path, cmd.HandleFunc)
	}
	registerBookListeners(r, b)
	registerAdminListeners(r, b)
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
	creates = append(creates, bookCommands...)
	creates = append(creates, adminCommands...)
	slog.Info("syncing commands...", slog.String("commands", fmt.Sprintf("%+v", creates)))
	err := handler.SyncCommands(b.client, creates, make([]snowflake.ID, 0))
	return err
}

// Logger is a middleware that logs the interaction and its variables.
var loggingMiddleware handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		slog.InfoContext(event.Ctx,
			"handling interaction",
			slog.String("interaction", fmt.Sprintf("%+v", event.Interaction)),
			slog.Any("vars", event.Vars),
		)
		return next(event)
	}
}
