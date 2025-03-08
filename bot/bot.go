package bot

import (
	"cubebot/bot/buttons"
	"cubebot/bot/commands"
	"cubebot/internal/db"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

var (
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func Run(db *sql.DB) {
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// open session first
	commandList := commands.GetCommands()
	discord.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		registerCommandsForGuild(s, g.Guild.ID, commandList)
	})
	err = discord.Open()
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	defer discord.Close() // close session, after function termination

	//slash commands
	commandHandlers := commands.GetCommandHandlers()
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Handle different types of interactions
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			// Create wrapper to handle user in database before processing command
			handleInteractionWithUserTracking(s, i, commandHandlers)
		case discordgo.InteractionMessageComponent:
			// Handle button clicks
			handleComponentInteraction(s, i)
		}
	})

	// add a event handler
	discord.AddHandler(newMessage)

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if *RemoveCommands {
		for _, guild := range discord.State.Guilds {
			removeCommandsForGuild(discord, guild.ID, commandList)
		}
	}

}

// handleInteractionWithUserTracking ensures the user is saved in the database
// before handling the command interaction
func handleInteractionWithUserTracking(s *discordgo.Session, i *discordgo.InteractionCreate,
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
	}

	// Get the user from the interaction
	user := i.User
	if user == nil && i.Member != nil {
		user = i.Member.User
	}

	if user != nil {
		// Insert or update user in database
		err := db.UpsertUser(user, guild)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			// Continue processing even if DB operation fails
		}
	}

	// After tracking the user, call the original handler
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	}
}

// handleComponentInteraction handles button clicks and other component interactions
func handleComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Track the user in the database first
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
	}

	// Get the user from the interaction
	user := i.User
	if user == nil && i.Member != nil {
		user = i.Member.User
	}

	if user != nil {
		// Insert or update user in database
		err := db.UpsertUser(user, guild)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			// Continue processing even if DB operation fails
		}
	}

	// After tracking the user, route to the appropriate button handler
	buttons.HandleButton(s, i)
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// TODO: Not sure if this will do anything ever

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

func removeCommandsForGuild(s *discordgo.Session, guildID string, commands []*discordgo.ApplicationCommand) {
	log.Println("Removing commands for guild:", guildID)

	// First, get all registered commands for this guild
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Printf("Warning: Failed to fetch commands for guild %s: %v", guildID, err)
		return
	}

	// Delete each command by name
	for _, cmd := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			// Log error but don't panic
			log.Printf("Warning: Cannot delete '%v' command in guild %s: %v", cmd.Name, guildID, err)
		} else {
			log.Printf("Command removed: %s", cmd.Name)
		}
	}
}
