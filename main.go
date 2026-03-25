package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/hazzardr/sbtbot/cmd/cli"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	slog.SetDefault(slog.New(logger))

	sbtb := cli.App{}
	ctx := kong.Parse(&sbtb,
		kong.Name("sbtbot"),
		kong.Description("Smart Bitches Trashy Bot CLI."),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	if err != nil {
		slog.Error("failed to start sbtbot", "error", err)
	}
}
