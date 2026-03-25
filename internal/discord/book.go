package discord

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var bookClubModal = discord.SlashCommandCreate{
	Name:        "submit",
	Description: "Submit a new book club idea!",
	IntegrationTypes: []discord.ApplicationIntegrationType{
		discord.ApplicationIntegrationTypeGuildInstall,
	},
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
}

func registerBookListeners(h *handler.Mux, b *Bot) {
	h.SlashCommand("/submit", b.launchSubmitModal)
}

var bookCommands = []discord.ApplicationCommandCreate{
	bookClubModal,
}

func (b *Bot) launchSubmitModal(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	genres, err := b.db.ListGenres(e.Ctx)
	if err != nil {
		return err
	}
	genreSelections := make([]discord.StringSelectMenuOption, len(genres))
	for i, g := range genres {
		genreSelections[i] = discord.NewStringSelectMenuOption(g.Name, strings.ToLower(g.Name))
	}
	genreSelect := discord.NewStringSelectMenu("genre-select", "Genre...", genreSelections...).
		WithMinValues(1).
		WithMaxValues(1)

	return e.Modal(discord.NewModalCreate("book-idea-modal", "Submit a new book club idea!", []discord.LayoutComponent{
		discord.NewLabel("Name", discord.NewShortTextInput("book-name")),
		discord.NewLabel("Author", discord.NewShortTextInput("author-name")),
		discord.NewLabel("Genre", genreSelect),
		discord.NewLabel("Theme", discord.NewStringSelectMenu("theme-select", "Theme...",
			discord.NewStringSelectMenuOption("Theme 1", "theme-1"),
			discord.NewStringSelectMenuOption("Theme 2", "theme-2"),
			discord.NewStringSelectMenuOption("Theme 3", "theme-3"),
		).WithMinValues(1).WithMaxValues(1)),
	}))
}
