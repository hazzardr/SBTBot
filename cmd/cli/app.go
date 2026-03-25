package cli

type App struct {
	Bot BotCommand `cmd:"" help:"Perform discord actions."`
}
