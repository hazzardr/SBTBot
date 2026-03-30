package discord

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var themeCommands = []discord.ApplicationCommandCreate{
	themeCRUD,
}

var themeCRUD = discord.SlashCommandCreate{
	Name:        "themes",
	Description: "Manage themes available to users",
	IntegrationTypes: []discord.ApplicationIntegrationType{
		discord.ApplicationIntegrationTypeGuildInstall,
	},
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "operation",
			Description: "What to do with the themes",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "List", Value: "list"},
				{Name: "Add", Value: "add"},
				{Name: "Delete", Value: "delete"},
			},
		},
		discord.ApplicationCommandOptionString{
			Name:        "name",
			Description: "Name of the theme. Not required for List",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "description",
			Description: "Description of the theme. Not required for List",
			Required:    false,
		},
	},
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionAdministrator),
}

func registerThemeListeners(h *handler.Mux, b *Bot) {
	h.SlashCommand("/themes", b.handleThemes)
}

func (b *Bot) handleThemes(d discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	operation := d.String("operation")
	name := d.String("name")
	description := d.String("description")
	if operation == "" {
		return errors.New("category option not found! this is a bug")
	}
	if operation == "list" {
		responseContent := ""
		themes, err := b.db.ListThemes(e.Ctx)
		if err != nil {
			slog.Error("error listing themes", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error getting themes, this is a bug"),
			)
		}
		if len(themes) == 0 {
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("No genres found"),
			)
		}
		for _, t := range themes {
			responseContent += "* " + t.Name + "\n"
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
		err := b.db.AddTheme(e.Ctx, name, description)
		if err != nil {
			slog.Error("error adding theme", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error adding theme, this is a bug"),
			)
		}
		return e.CreateMessage(discord.NewMessageCreate().
			WithEphemeral(true).
			WithContentf("✅ Created Theme: %s", name),
		)
	} else if operation == "delete" {
		if name == "" {
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("❌ Must provide a name!"),
			)
		}
		err := b.db.RemoveTheme(e.Ctx, name)
		if err != nil {
			slog.Error("error removing theme", err)
			return e.CreateMessage(discord.NewMessageCreate().
				WithEphemeral(true).
				WithContent("Error removing theme, this is a bug"),
			)
		}
		return e.CreateMessage(discord.NewMessageCreate().
			WithEphemeral(true).
			WithContentf("✅ Removed Theme: %s", name),
		)
	} else {
		return errors.New("unknown operation: " + operation)
	}
}
