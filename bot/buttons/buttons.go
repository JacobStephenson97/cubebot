package buttons

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ButtonHandler is a function that handles a button interaction
type ButtonHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// GetButtonHandlers returns a map of button handlers
// The key is the prefix of the CustomID, e.g. "join_queue"
func GetButtonHandlers() map[string]ButtonHandler {
	return map[string]ButtonHandler{
		"join_queue":  HandleJoinDraft,
		"leave_queue": HandleLeaveQueue,
	}
}

// HandleButton routes the button interaction to the appropriate handler based on the CustomID
func HandleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract the button ID and parameters
	customID := i.MessageComponentData().CustomID
	parts := strings.SplitN(customID, "_", 3) // Split into max 3 parts: action_type_params

	if len(parts) < 2 {
		// Invalid button ID format
		return
	}

	// The first two parts form the button type (e.g., "join_queue")
	buttonType := parts[0] + "_" + parts[1]

	// Find and call the appropriate handler
	if handler, ok := GetButtonHandlers()[buttonType]; ok {
		handler(s, i)
	}
}
