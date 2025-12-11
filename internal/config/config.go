package config

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	IntervalsBaseUrl   string
	IntervalsAthleteId string
	IntervalsApiKey    string
	McpAuthToken       string
	Port               string
}

func Load() Config {
	// notice github.com/joho/godotenv/autoload

	config := Config{
		IntervalsBaseUrl:   os.Getenv("INTERVALS_API_BASE_URL"),
		IntervalsAthleteId: os.Getenv("INTERVALS_ATHLETE_ID"),
		IntervalsApiKey:    os.Getenv("INTERVALS_API_KEY"),
		McpAuthToken:       os.Getenv("MCP_AUTH_TOKEN"),
		Port:               os.Getenv("PORT"),
	}

	if config.IntervalsBaseUrl == "" {
		config.IntervalsBaseUrl = "https://intervals.icu"
	}
	if config.Port == "" {
		config.Port = "8000"
	}

	if len(config.IntervalsApiKey) < 4 {
		log.Panic("Configuring INTERVALS_API_KEY required.")
	}
	if len(config.IntervalsAthleteId) < 4 {
		log.Panic("Configuring INTERVALS_ATHLETE_ID required.")
	}
	if len(config.McpAuthToken) < 20 {
		token, err := generatePassword()
		if err != nil {
			log.Panic("Configuring MCP_AUTH_TOKEN required.")
		}
		config.McpAuthToken = token
		log.Println("MCP_AUTH_TOKEN was shorter than 20 chars. Randomising.")
		log.Printf("MCP_AUTH_TOKEN=%s", config.McpAuthToken)
	}

	fmt.Printf("Config: %s=%v\n", "IntervalsBaseUrl", config.IntervalsBaseUrl)
	fmt.Printf("Config: %s=%v...\n", "IntervalsAthleteId", config.IntervalsAthleteId[0:3])
	fmt.Printf("Config: %s=%v...\n", "IntervalsApiKey", config.IntervalsApiKey[0:3])
	fmt.Printf("Config: %s=%v...\n", "McpAuthToken", config.McpAuthToken[0:3])
	fmt.Printf("Config: %s=%v\n", "Port", config.Port)

	return config
}

// Define the character set for the password
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const passwordLength = 32

func generatePassword() (string, error) {
    // Convert the character set string to a slice of runes (for easy indexing)
    charsetRunes := []rune(charset)
    
    // Get the length of the character set
    charsetLen := big.NewInt(int64(len(charsetRunes)))
    
    // Create a byte slice to hold the generated password characters
    password := make([]rune, passwordLength)

    // Loop until all characters in the password have been selected
    for i := range passwordLength {
        // Generate a cryptographically secure random number in the range [0, charsetLen-1]
        randomIndex, err := rand.Int(rand.Reader, charsetLen)
        if err != nil {
            return "", fmt.Errorf("error generating random number: %w", err)
        }
        
        // Use the random number as an index to select a character from the set
        password[i] = charsetRunes[randomIndex.Int64()]
    }

    // Convert the rune slice back into a string and return it
    return string(password), nil
}