package buttons

import (
	"cubebot/bot/commands"
	"cubebot/internal/db"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleLeaveQueue handles the leave_queue button interaction
func HandleLeaveQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract the session ID from the custom ID
	// Format: leave_queue_sessionID
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, "_")

	if len(parts) < 3 {
		log.Printf("Invalid leave_queue button ID format: %s", customID)
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

	// Check if the user is in the draft queue
	var count int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM draft_participants WHERE session_id = ? AND user_id = ?", sessionID, user.ID).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user is in draft queue: %v", err)
		respondWithError(s, i, "Error leaving queue. Please contact an administrator.")
		return
	}

	if count == 0 {
		respondWithError(s, i, "You're not in the draft queue!")
		return
	}

	// Remove the user from the draft queue
	_, err = db.GetDB().Exec("DELETE FROM draft_participants WHERE session_id = ? AND user_id = ?", sessionID, user.ID)
	if err != nil {
		log.Printf("Error removing user from draft queue: %v", err)
		respondWithError(s, i, "Error leaving queue. Please contact an administrator.")
		return
	}

	// Update the session embed to reflect the changes
	err = commands.UpdateSessionEmbed(s, sessionID)
	if err != nil {
		log.Printf("Error updating session embed: %v", err)
		respondWithError(s, i, "Error leaving queue. Please contact an administrator.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You've left the draft queue!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error sending followup message: %v", err)
	}
}
