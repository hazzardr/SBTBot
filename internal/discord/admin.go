package discord

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var adminCommands = []discord.ApplicationCommandCreate{
	adminModal,
}

var adminModal = discord.SlashCommandCreate{
	Name:        "manage",
	Description: "Manage themes, genres, and submissions",
	IntegrationTypes: []discord.ApplicationIntegrationType{
		discord.ApplicationIntegrationTypeGuildInstall,
	},
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "category",
			Description: "The category to manage",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "Themes", Value: "themes"},
				{Name: "Genres", Value: "genres"},
			},
		},
	},
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionAdministrator),
}

func registerAdminListeners(h *handler.Mux, b *Bot) {
	h.SlashCommand("/manage", b.launchAdminModal)
}

func (b *Bot) launchAdminModal(d discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	slog.Info("launching modal", slog.String("data", fmt.Sprintf("%+v", d)))
	modalType, ok := d.Options["category"]
	if !ok {
		return errors.New("category option not found! this is a bug")
	}
	am, err := b.adminModal(e.Ctx, modalType.String())
	if err != nil {
		return fmt.Errorf("failed to render admin modal: %w", err)
	}
	return e.Modal(am)
}

func (b *Bot) adminModal(ctx context.Context, modalType string) (discord.ModalCreate, error) {
	opts := make([]discord.StringSelectMenuOption, 0)
	if modalType == "themes" {
		themes, err := b.db.ListThemes(ctx)
		if err != nil {
			return discord.ModalCreate{}, err
		}
		for _, t := range themes {
			opts = append(opts, discord.NewStringSelectMenuOption(t.Name, t.Name))
		}
	}
	return discord.NewModalCreate("admin-modal", fmt.Sprintf("Manage %s", modalType), []discord.LayoutComponent{
		discord.NewLabel("Delete", discord.NewStringSelectMenu("theme-delete", "",
			opts...,
		)),
	}), nil
}
