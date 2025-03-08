package commands

import (
	"cubebot/internal/db"
	"log"

	"github.com/bwmarrin/discordgo"
)

func RemoveCubeCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "remove-cube",
		Description: "Remove a cube from the database",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "cube",
				Description: "The cube to remove",
				Required:    true,
				Choices:     GetDraftFormatChoices(),
			},
		},
	}
}

func HandleRemoveCube(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	// Get the cube option value
	var cubeName string
	for _, option := range options {
		if option.Name == "cube" {
			cubeName = option.StringValue()
			break
		}
	}

	exists, err := cubeExists(cubeName)
	if err != nil {
		log.Printf("Error checking if cube exists: %v", err)
		return
	}
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cube not found",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	_, err = db.GetDB().Exec("DELETE FROM draft_formats WHERE name = ?", cubeName)
	if err != nil {
		log.Printf("Error removing cube from database: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error removing cube from database",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cube removed from database",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	registerCommandsForGuild(s, i.GuildID, GetCommands())
}
func registerCommandsForGuild(s *discordgo.Session, guildID string, commandList []*discordgo.ApplicationCommand) {
	log.Println("Registering commands for guild: ", guildID)
	for _, v := range commandList {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		log.Println("Command registered: ", cmd.Name)
	}
}
