package discord

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func getLaunchSubmitModalCommand(b *Bot) SlashCommand {
	return SlashCommand{
		Metadata: discord.SlashCommandCreate{
			Name:        "submit",
			Description: "Launches a pop-up to submit ideas",
			IntegrationTypes: []discord.ApplicationIntegrationType{
				discord.ApplicationIntegrationTypeGuildInstall,
			},
			Contexts: []discord.InteractionContextType{
				discord.InteractionContextTypeGuild,
			},
		},
		Path:       "/submit",
		HandleFunc: b.launchSubmitModal,
	}
}

func (b *Bot) launchSubmitModal(d discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	// b.getGenres
	return e.Modal(discord.NewModalCreate("book-idea-modal", "Submit a new book club idea!", []discord.LayoutComponent{
		discord.NewLabel("Name", discord.NewShortTextInput("book-name")),
		discord.NewLabel("Author", discord.NewShortTextInput("author-name")),
		discord.NewLabel("Genre", discord.NewStringSelectMenu("genre-select", "Genre...",
			discord.NewStringSelectMenuOption("Fantasy", "fantasy"),
			discord.NewStringSelectMenuOption("Historical", "historical"),
			discord.NewStringSelectMenuOption("Contemporary", "contemporary"),
		).WithMinValues(1).WithMaxValues(1)),
		discord.NewLabel("Theme", discord.NewStringSelectMenu("theme-select", "Theme...",
			discord.NewStringSelectMenuOption("Theme 1", "theme-1"),
			discord.NewStringSelectMenuOption("Theme 2", "theme-2"),
			discord.NewStringSelectMenuOption("Theme 3", "theme-3"),
		).WithMinValues(1).WithMaxValues(1)),
	}))
}
