package commands

import (
	"github.com/bwmarrin/discordgo"
)

// GetCommands returns all the commands for the bot
func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		StartDraftCommand(),
		AddCubeCommand(),
		RemoveCubeCommand(),
	}
}

// GetCommandHandlers returns a map of command handlers
func GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start-draft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			HandleStartDraft(s, i)
		},
		"add-cube": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			HandleAddCube(s, i)
		},
		"remove-cube": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			HandleRemoveCube(s, i)
		},
	}
}
