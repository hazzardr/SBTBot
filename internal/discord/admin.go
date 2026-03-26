package discord

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (b *Bot) launchAdminModal(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.Modal(discord.NewModalCreate("admin-modal", "Admin", []discord.LayoutComponent{
		discord.NewTextDisplay("Add and remove choices for genres and themes."),
	}))
}
