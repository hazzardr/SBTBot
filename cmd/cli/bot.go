package cli

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hazzardr/sbtbot/internal/discord"
)

type BotCommand struct {
	Serve ServeCommand `cmd:"" help:"Start the Discord Bot"`
	Sync  SyncCommand  `cmd:"" help:"Sync the Discord Application Commands. Necessary when adding a new slash command."`
}

type ServeCommand struct {
	Token string `required:"true" help:"Discord bot token"`
	DSN   string `required:"true" default:"data/sbtb.db" help:"database path"`
}
type SyncCommand struct {
	Token string `required:"true" help:"Discord bot token"`
	DSN   string `required:"true" default:"data/sbtb.db" help:"database path"`
}

func (s *ServeCommand) Run() error {
	b, err := discord.NewBot(s.Token, s.DSN)
	if err != nil {
		return fmt.Errorf("failed to create discord bot: %w", err)
	}
	err = b.Start()
	if err != nil {
		return fmt.Errorf("failed to start the discord bot: %w", err)
	}
	slog.Info("discord connection successful")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig
	return nil
}

func (s *SyncCommand) Run() error {
	b, err := discord.NewBot(s.Token, s.DSN)
	if err != nil {
		return fmt.Errorf("failed to create discord bot: %w", err)
	}
	err = b.SyncCommands()
	if err != nil {
		return fmt.Errorf("failed to sync commands: %w", err)
	}
	slog.Info("discord command sync successful")
	return nil
}
