package bot

import (
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
var database *sql.DB // Add a package-level variable to store the DB connection

var (
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run(db *sql.DB) {
	database = db // Store the DB connection
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// open session first
	err = discord.Open()
	checkNilErr(err)
	defer discord.Close() // close session, after function termination

	//slash commands
	var commands = []*discordgo.ApplicationCommand{
		{
			Name: "start-draft",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Start a new draft",
		},
	}
	var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start-draft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Draft started",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Create wrapper to handle user in database before processing command
		handleInteractionWithUserTracking(s, i, commandHandlers)
	})

	//add a event handler for guild join
	discord.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		registerCommandsForGuild(s, g.Guild.ID, commands)
	})

	// Register commands after opening the session

	// add a event handler
	discord.AddHandler(newMessage)

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if *RemoveCommands {
		for _, guild := range discord.State.Guilds {
			removeCommandsForGuild(discord, guild.ID, commands)
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
		err := db.UpsertUser(database, user, guild)
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

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// TODO: Not sure if this will do anything ever

}

func registerCommandsForGuild(s *discordgo.Session, guildID string, commands []*discordgo.ApplicationCommand) {
	log.Println("Registering commands for guild: ", guildID)
	for _, v := range commands {
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
