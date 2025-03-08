package buttons

import (
	"cubebot/bot/commands"
	"cubebot/internal/db"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleJoinDraft handles the join_queue button interaction
func HandleJoinDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract the session ID from the custom ID
	// Format: join_queue_sessionID
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, "_")

	if len(parts) < 3 {
		log.Printf("Invalid join_queue button ID format: %s", customID)
		respondWithError(s, i, "Invalid button format. Please contact an administrator.")
		return
	}

	sessionIDStr := parts[2]
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		log.Printf("Error parsing session ID from button: %v", err)
		respondWithError(s, i, "Invalid session ID. Please contact an administrator.")
		return
	}

	// Get the user who clicked the button
	user := i.User
	if user == nil && i.Member != nil {
		user = i.Member.User
	}

	// Acknowledge the interaction immediately
	// err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Flags: discordgo.MessageFlagsEphemeral, // Make the response visible only to the user who clicked
	// 	},
	// })
	// if err != nil {
	// 	log.Printf("Error responding to interaction: %v", err)
	// 	return
	// }

	//check if the user is already in the draft queue
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM draft_participants WHERE session_id = ? AND user_id = ?", sessionID, user.ID).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user is in draft queue: %v", err)
		respondWithError(s, i, "Error joining queue. Please contact an administrator.")
		return
	}

	if count > 0 {
		respondWithError(s, i, "You're already in the draft queue!")
		return
	}
	// Update the draft session to include the user
	err = db.AddParticipantToDraftSession(sessionID, user.ID)
	if err != nil {
		log.Printf("Error updating draft session: %v", err)
		respondWithError(s, i, "Error joining queue. Please contact an administrator.")
		return
	}

	err = commands.UpdateSessionEmbed(s, sessionID)
	if err != nil {
		log.Printf("Error getting updated session embed: %v", err)
		respondWithError(s, i, "Error joining queue. Please contact an administrator.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You've joined the draft queue!",
			Flags:   discordgo.MessageFlagsEphemeral, // Make the error visible only to the user who clicked
		},
	})
	if err != nil {
		log.Printf("Error sending followup message: %v", err)
	}
}

// Helper function to respond with an error message
func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, errorMsg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + errorMsg,
			Flags:   discordgo.MessageFlagsEphemeral, // Make the error visible only to the user who clicked
		},
	})
	if err != nil {
		log.Printf("Error responding with error message: %v", err)
	}
}
