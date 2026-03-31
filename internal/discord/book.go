package discord

import (
	"context"
	"fmt"
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
		discord.NewLabel("Name", discord.NewShortTextInput("title")),
		discord.NewLabel("Author", discord.NewShortTextInput("author")),
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
	genreSelect = discord.NewStringSelectMenu("genre", "Genre...", genreSelections...).
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
	title, ok := e.Data.TextInput("title")
	if !ok {
		return e.CreateMessage(discord.NewMessageCreate().
			WithContent("Please enter a title!").
			WithEphemeral(true),
		)
	}
	author, ok := e.Data.TextInput("author")
	if !ok {
		return e.CreateMessage(discord.NewMessageCreate().
			WithContent("Please enter an author!").
			WithEphemeral(true),
		)
	}
	genre, _ := e.Data.StringSelectMenu("genre")
	var g string
	if len(genre.Values) != 1 {
		g = ""
	} else {
		g = genre.Values[0]
	}
	theme, _ := e.Data.StringSelectMenu("theme")
	var t string
	if len(theme.Values) != 1 {
		t = ""
	} else {
		t = theme.Values[0]
	}
	submitter := e.User().Username
	bi, err := b.db.AddBookIdea(
		e.Ctx,
		title.Value,
		author.Value,
		g,
		t,
		submitter,
	)
	if err != nil {
		return err
	}
	return e.CreateMessage(discord.NewMessageCreate().
		WithContent(fmt.Sprintf("Submitted book: %s. Thanks!", bi.Title)).
		WithEphemeral(true),
	)
}
