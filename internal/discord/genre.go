package discord

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var genreCommands = []discord.ApplicationCommandCreate{
	genreCRUD,
}

var genreCRUD = discord.SlashCommandCreate{
	Name:        "genres",
	Description: "Manage genres available to users",
	IntegrationTypes: []discord.ApplicationIntegrationType{
		discord.ApplicationIntegrationTypeGuildInstall,
	},
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "operation",
			Description: "What to do with the genres",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "List", Value: "list"},
				{Name: "Add", Value: "add"},
				{Name: "Delete", Value: "delete"},
			},
		},
		discord.ApplicationCommandOptionString{
			Name:        "name",
			Description: "Name of the genre. Not required for List",
			Required:    false,
		},
	},
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionAdministrator),
}

func registerGenreListeners(h *handler.Mux, b *Bot) {
	h.SlashCommand("/genres", b.handleGenres)
}

func (b *Bot) handleGenres(d discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	operation := d.String("operation")
	name := d.String("name")
	if operation == "" {
		return errors.New("category option not found! this is a bug")
	}
	if operation == "list" {
		responseContent := ""
		genres, err := b.db.ListGenres(e.Ctx)
		if err != nil {
			slog.Error("error listing genres", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error getting genres, this is a bug"),
			)
		}
		if len(genres) == 0 {
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("No genres found"),
			)
		}
		for _, genre := range genres {
			responseContent += "* " + genre.Name + "\n"
		}
		return e.CreateMessage(discord.NewMessageCreate().
			WithEphemeral(true).
			WithContent(responseContent),
		)
	} else if operation == "add" {
		if name == "" {
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("❌ Must provide a name!"),
			)
		}
		err := b.db.AddGenre(e.Ctx, name)
		if err != nil {
			slog.Error("error adding genre", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error adding genre, this is a bug"),
			)
		}
		return e.CreateMessage(discord.NewMessageCreate().
			WithEphemeral(true).
			WithContentf("✅ Created Genre: %s", name),
		)
	} else if operation == "delete" {
		if name == "" {
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("❌ Must provide a name!"),
			)
		}
		err := b.db.RemoveGenre(e.Ctx, name)
		if err != nil {
			slog.Error("error removing genre", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error removing genre, this is a bug"),
			)
		}
		return e.CreateMessage(discord.NewMessageCreate().
			WithEphemeral(true).
			WithContentf("✅ Removed Genre: %s", name),
		)
	} else {
		return errors.New("unknown operation: " + operation)
	}
}
