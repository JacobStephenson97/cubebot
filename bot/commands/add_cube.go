package commands

import (
	"cubebot/internal/db"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddCubeCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "add-cube",
		Description: "Add a cube to the database",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "cube",
				Description: "The cube to add",
				Required:    true,
			},
		},
	}
}

func HandleAddCube(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	// Get the cube option value
	var cubeName string
	for _, option := range options {
		if option.Name == "cube" {
			cubeName = option.StringValue()
			break
		}
	}

	resp, err := http.Get("https://cubecobra.com/cube/list/" + cubeName)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	if strings.Contains(string(body), "404: Page not found") {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cube not found",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	cubeURL := "https://cubecobra.com/cube/list/" + cubeName

	exists, err := cubeExists(cubeName)
	if err != nil {
		log.Printf("Error checking if cube exists: %v", err)
		return
	}
	if exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cube already exists",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	_, err = db.GetDB().Exec("INSERT INTO draft_formats (name, cube_url) VALUES (?, ?)", cubeName, cubeURL)

	if err != nil {
		log.Printf("Error adding cube to database: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error adding cube to database",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cube added to database",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	registerCommandsForGuild(s, i.GuildID, GetCommands())
}

func cubeExists(cubeName string) (bool, error) {
	var exists bool
	err := db.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM draft_formats WHERE name = ?)", cubeName).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if cube exists: %v", err)
		return false, err
	}

	return exists, nil
}
