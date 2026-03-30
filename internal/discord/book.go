package discord

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/hazzardr/sbtbot/internal/sbtb"
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
	h.Modal("/submit-idea", b.submitBookIdea)
}

var bookCommands = []discord.ApplicationCommandCreate{
	bookClubModal,
}

func (b *Bot) launchSubmitModal(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	genreSelect, err := genreDropDown(e.Ctx, b.db)
	if err != nil {
		return fmt.Errorf("failed to render genre component: %w", err)
	}
	themeSelect, err := themeDropDown(e.Ctx, b.db)
	if err != nil {
		return fmt.Errorf("failed to render theme component: %w", err)
	}

	return e.Modal(discord.NewModalCreate("/submit-idea", "Submit a new book club idea!", []discord.LayoutComponent{
		discord.NewLabel("Name", discord.NewShortTextInput("book-name")),
		discord.NewLabel("Author", discord.NewShortTextInput("author-name")),
		genreSelect,
		themeSelect,
	}))
}

func genreDropDown(ctx context.Context, db *sbtb.DB) (discord.LayoutComponent, error) {
	var genreSelect discord.LabelSubComponent
	genres, err := db.ListGenres(ctx)
	if err != nil {
		return nil, err
	}
	if len(genres) == 0 {
		return discord.NewTextDisplay("**No genres** created yet!"), nil
	}
	genreSelections := make([]discord.StringSelectMenuOption, len(genres))
	for i, g := range genres {
		opt := discord.NewStringSelectMenuOption(g.Name, strings.ToLower(g.Name))
		genreSelections[i] = opt
	}
	genreSelect = discord.NewStringSelectMenu("genre-select", "Genre...", genreSelections...).
		WithMinValues(1).
		WithMaxValues(1).
		WithRequired(true)
	return discord.NewLabel("Genre", genreSelect), nil
}

func themeDropDown(ctx context.Context, db *sbtb.DB) (discord.LayoutComponent, error) {
	var themeSelect discord.LabelSubComponent
	themes, err := db.ListThemes(ctx)
	if err != nil {
		return nil, err
	}
	if len(themes) == 0 {
		return discord.NewTextDisplay("**No themes** created yet!"), nil
	}
	themeSelections := make([]discord.StringSelectMenuOption, len(themes))
	for i, t := range themes {
		opt := discord.NewStringSelectMenuOption(t.Name, strings.ToLower(t.Name))
		if t.Description.Valid {
			opt.Description = t.Description.String
		}
		themeSelections[i] = opt
	}
	themeSelect = discord.NewStringSelectMenu("theme-select", "Theme...", themeSelections...).
		WithMinValues(1).
		WithMaxValues(1)
	return discord.NewLabel("Theme", themeSelect), nil
}

func (b *Bot) submitBookIdea(e *handler.ModalEvent) error {
	slog.InfoContext(e.Ctx, "Submitting book idea")
	return nil
}
