package commands

import (
	"cubebot/bot/utils"
	"cubebot/internal/db"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GetDraftFormatChoices returns the choices for the draft format command
func GetDraftFormatChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}

	// Query the database for draft formats
	rows, err := db.GetDB().Query("SELECT name FROM draft_formats")
	if err != nil {
		log.Printf("Error querying draft formats: %v", err)
		return choices
	}
	defer rows.Close()

	// Iterate through the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Error scanning draft format row: %v", err)
			continue
		}

		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: name,
		})
	}

	return choices
}

// StartDraftCommand returns the command definition for the start-draft command
func StartDraftCommand() *discordgo.ApplicationCommand {

	choices := GetDraftFormatChoices()
	return &discordgo.ApplicationCommand{
		Name:        "start-draft",
		Description: "Start a new draft",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "cube",
				Description: "The cube to draft from",
				Required:    true,
				Choices:     choices,
			},
		},
	}
}

// HandleStartDraft handles the start-draft command interaction
func HandleStartDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Send a deferred response first
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting draft...",
		},
	})
	if err != nil {
		log.Printf("Error sending deferred response: %v", err)
		return
	}

	// Get the cube option value
	var cubeName string
	for _, option := range options {
		if option.Name == "cube" {
			cubeName = option.StringValue()
			break
		}
	}

	// Query the database for the selected cube
	var cubeURL string
	var formatID int
	err = db.GetDB().QueryRow("SELECT id, cube_url FROM draft_formats WHERE name = ?", cubeName).Scan(&formatID, &cubeURL)
	if err != nil {
		log.Printf("Error querying cube URL: %v", err)
		errorMsg := "Error: Could not find cube information for " + cubeName
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errorMsg,
		})
		if err != nil {
			log.Printf("Error editing interaction response: %v", err)
		}
		return
	}

	//Create the draft session
	randomCode := utils.GenerateRandomString(8)
	externDraftURL := "https://draftmancer.com/?session=" + randomCode

	// Use Exec and LastInsertId to get the ID of the newly created session
	var sessionID int64
	result, err := db.GetDB().Exec("INSERT INTO draft_sessions (format_id, guild_id, created_by_user_id, status, external_draft_url, channel_id) VALUES (?, ?, ?, ?, ?, ?)",
		formatID, i.GuildID, i.Member.User.ID, "queue", externDraftURL, i.ChannelID)

	if err != nil {
		log.Printf("Error creating draft session: %v", err)
		errorMsg := "Error: Could not create draft session"
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errorMsg,
		})
		if err != nil {
			log.Printf("Error editing interaction response: %v", err)
		}
		return
	}

	// Get the last inserted ID
	sessionID, err = result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		// Continue anyway, we just won't have the session ID in the embed
	}

	// Create the embed for the response
	embed := &discordgo.MessageEmbed{
		Title:       "Draft Started: " + cubeName,
		Description: "A new draft session has been created. Click the link below to join!",
		URL:         cubeURL,
		Color:       0x00BFFF, // Deep Sky Blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Cube",
				Value:  cubeName,
				Inline: true,
			},
			{
				Name:   "Cube URL",
				Value:  "[View Cube](" + cubeURL + ")",
				Inline: true,
			},
			{
				Name:   "Draft URL",
				Value:  "[Join Draft](" + externDraftURL + ")",
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  "Waiting for players",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Created by " + i.Member.User.Username,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Create the components (buttons)
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Join Queue",
					Style:    discordgo.PrimaryButton,
					CustomID: "join_queue_" + fmt.Sprintf("%d", sessionID),
					Disabled: false,
				},
				discordgo.Button{
					Label:    "Leave Queue",
					Style:    discordgo.DangerButton,
					CustomID: "leave_queue_" + fmt.Sprintf("%d", sessionID),
					Disabled: false,
				},
				discordgo.Button{
					Label:    "View Cube",
					Style:    discordgo.LinkButton,
					URL:      cubeURL,
					Disabled: false,
				},
			},
		},
	}

	// Edit the deferred response with our embed and components
	message, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
	if err != nil {
		log.Printf("Error editing interaction response with embed: %v", err)
	}

	// Update the draft session with the message ID
	_, err = db.GetDB().Exec("UPDATE draft_sessions SET message_id = ? WHERE id = ?", message.ID, sessionID)
	if err != nil {
		log.Printf("Error updating draft session with message ID: %v", err)
	}
}

func UpdateSessionEmbed(s *discordgo.Session, sessionID int) error {

	var cubeName string
	var cubeURL string
	var externDraftURL string
	var channelID string
	var messageID string
	var status string
	var formatID int
	//get the session data from draftsessions table
	err := db.GetDB().QueryRow("SELECT format_id, external_draft_url, channel_id, message_id, status FROM draft_sessions WHERE id = ?", sessionID).Scan(&formatID, &externDraftURL, &channelID, &messageID, &status)
	if err != nil {
		log.Printf("Error querying draft session: %v", err)
		return err
	}
	//get the cube name from the draft_formats table
	err = db.GetDB().QueryRow("SELECT name, cube_url FROM draft_formats WHERE id = ?", formatID).Scan(&cubeName, &cubeURL)
	if err != nil {
		log.Printf("Error querying draft format: %v", err)
		return err
	}

	//get the participants from the draft_participants table
	participants, err := db.GetDB().Query("SELECT user_id FROM draft_participants WHERE session_id = ?", sessionID)
	if err != nil {
		log.Printf("Error querying draft participants: %v", err)
		return err
	}
	defer participants.Close()

	var participantList string
	for participants.Next() {
		var userID string
		err = participants.Scan(&userID)
		if err != nil {
			log.Printf("Error scanning draft participant: %v", err)
			continue
		}
		participantList += "<@" + userID + ">\n"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Draft Started: " + cubeName,
		Description: "A new draft session has been created. Click the link below to join!",
		URL:         cubeURL,
		Color:       0x00BFFF, // Deep Sky Blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Cube",
				Value:  cubeName,
				Inline: true,
			},
			{
				Name:   "Cube URL",
				Value:  "[View Cube](" + cubeURL + ")",
				Inline: true,
			},
			{
				Name:   "Draft URL",
				Value:  "[Join Draft](" + externDraftURL + ")",
				Inline: true,
			},
			{
				Name:   "Participants",
				Value:  participantList,
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	//update the message with the new embed
	_, err = s.ChannelMessageEditEmbed(channelID, messageID, embed)
	if err != nil {
		log.Printf("Error updating message: %v", err)
		return err
	}
	return nil
}
