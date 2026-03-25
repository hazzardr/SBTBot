package discord

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var pingCommand = SlashCommand{
	Metadata: discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Checks bot health",
	},
	Path:       "/ping",
	HandleFunc: HandlePing,
}

func HandlePing(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var gatewayPing string
	if e.Client().HasGateway() {
		gatewayPing = e.Client().Gateway.Latency().String()
	}
	return e.Respond(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("**Pong!**"),
			discord.NewTextDisplayf("Latency: %s", gatewayPing),
		).WithAccentColor(colorSuccess),
	))
}
