package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Connect() (*sql.DB, error) {
	// Open up our database connection
	db, err := sql.Open("mysql", os.Getenv("DB_CONN_STRING"))

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	// Bootstrap the database schema
	if err := CreateTables(db); err != nil {
		return nil, fmt.Errorf("failed to bootstrap database: %w", err)
	}

	database = db
	return db, nil
}

// UpsertUser inserts a user into the database if they don't exist,
// or updates their information if they do exist
func UpsertUser(user *discordgo.User, guild *discordgo.Guild) error {
	currentTime := time.Now()

	// First, make sure the guild exists in the database
	if guild.ID != "" {
		var guildExists bool
		err := GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM guilds WHERE id = ?)", guild.ID).Scan(&guildExists)
		if err != nil {
			return fmt.Errorf("error checking if guild exists: %w", err)
		}

		if !guildExists {
			// If we don't have guild name available, just use the ID as a placeholder
			// Typically you would fetch the guild info using the Discord API
			_, err = GetDB().Exec(
				"INSERT INTO guilds (id, guild_name) VALUES (?, ?)",
				guild.ID,
				guild.Name, // Placeholder name
			)
			if err != nil {
				return fmt.Errorf("error inserting guild: %w", err)
			}
		}
	}

	// Then handle the user
	var userExists bool
	err := GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", user.ID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	// Construct avatar URL if available
	var avatarURL string
	if user.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.gif", user.ID, user.Avatar)
	}

	if userExists {
		// Update existing user
		_, err = GetDB().Exec(
			"UPDATE users SET discord_username = ?, display_name = ?, avatar_url = ?, updated_at = ? WHERE id = ?",
			user.Username,
			user.GlobalName,
			avatarURL,
			currentTime,
			user.ID,
		)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}
	} else {
		// Insert new user
		_, err = GetDB().Exec(
			"INSERT INTO users (id, discord_username, display_name, avatar_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			user.ID,
			user.Username,
			user.GlobalName,
			avatarURL,
			currentTime,
			currentTime,
		)
		if err != nil {
			return fmt.Errorf("error inserting new user: %w", err)
		}
	}

	return nil
}

// CreateTables initializes the database by running all SQL schema files in the tables directory
func CreateTables(db *sql.DB) error {
	// Path to the tables directory
	tablesDir := "./internal/db/tables"

	files, err := os.ReadDir(tablesDir)
	if err != nil {
		return fmt.Errorf("error reading tables directory: %w", err)
	}

	fmt.Println("Initializing database schema...")

	// Process SQL files in alphabetical order
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(tablesDir, file.Name())

			// Read the SQL file
			sqlBytes, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading SQL file %s: %w", filePath, err)
			}

			sqlContent := string(sqlBytes)
			fmt.Printf("Executing SQL file: %s\n", file.Name())

			// Execute the SQL
			_, err = db.Exec(sqlContent)
			if err != nil {
				// Ignore duplicate index errors (MySQL error 1061)
				// This lets us run the scripts multiple times without error
				if strings.Contains(err.Error(), "Error 1061") || strings.Contains(err.Error(), "Duplicate key name") {
				} else {
					// Print the SQL that failed to help debug other errors
					fmt.Printf("Error executing SQL from file %s:\n%s\n", file.Name(), sqlContent)
					return fmt.Errorf("error executing SQL from file %s: %w", filePath, err)
				}
			} else {
				fmt.Printf("Successfully executed SQL from file: %s\n", file.Name())
			}
		}
	}

	fmt.Println("Database bootstrap completed successfully")
	return nil
}

func GetDB() *sql.DB {
	return database
}

func AddParticipantToDraftSession(sessionID int, userID string) error {
	_, err := GetDB().Exec("INSERT INTO draft_participants (session_id, user_id) VALUES (?, ?)", sessionID, userID)
	return err
}

func RemoveParticipantFromDraftSession(sessionID int, userID string) error {
	_, err := GetDB().Exec("DELETE FROM draft_participants WHERE session_id = ? AND user_id = ?", sessionID, userID)
	return err
}

func GetDraftSessionParticipants(sessionID int) ([]string, error) {
	rows, err := GetDB().Query("SELECT user_id FROM draft_participants WHERE session_id = ?", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	participants := []string{}
	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}

	return participants, nil
}
