package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/St0iK/go-quote-bot/dao"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Config holds the application configuration including Twitter credentials
type Config struct {
	TwitterCreds TwitterCredentials
}

// TwitterCredentials stores Twitter API authentication tokens
type TwitterCredentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// TwitterClient wraps the Twitter client and provides application-specific functionality
type TwitterClient struct {
	client *twitter.Client
}

// NewConfig creates a new Config instance from environment variables
func NewConfig() (*Config, error) {
	creds := TwitterCredentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	// Validate credentials
	if creds.AccessToken == "" || creds.AccessTokenSecret == "" ||
		creds.ConsumerKey == "" || creds.ConsumerSecret == "" {
		return nil, fmt.Errorf("missing required Twitter credentials in environment variables")
	}

	return &Config{
		TwitterCreds: creds,
	}, nil
}

// NewTwitterClient creates and authenticates a new Twitter client
func NewTwitterClient(creds TwitterCredentials) (*TwitterClient, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Twitter credentials: %w", err)
	}

	log.Printf("Authenticated as Twitter user: @%s", user.ScreenName)
	return &TwitterClient{
		client: client,
	}, nil
}

// PostQuote posts a quote to Twitter
func (t *TwitterClient) PostQuote(quote *dao.Quote) error {
	if quote == nil {
		return fmt.Errorf("cannot post nil quote")
	}

	tweetText := formatTweet(quote)
	tweet, _, err := t.client.Statuses.Update(tweetText, nil)
	if err != nil {
		return fmt.Errorf("failed to post tweet: %w", err)
	}

	log.Printf("Successfully posted tweet with ID: %d", tweet.ID)
	return nil
}

// formatTweet creates a properly formatted tweet from a quote
func formatTweet(quote *dao.Quote) string {
	return fmt.Sprintf("%s - %s #quotes #quote #inspiration",
		quote.QuoteText,
		quote.Author)
}

func main() {
	log.Println("Starting Twitter Quote Bot v0.02")

	// Initialize configuration
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize database connection
	if err := dao.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get random quote
	quote, err := dao.GetRandomQuote()
	if err != nil {
		log.Fatalf("Failed to get random quote: %v", err)
	}

	// Initialize Twitter client
	client, err := NewTwitterClient(config.TwitterCreds)
	if err != nil {
		log.Fatalf("Failed to initialize Twitter client: %v", err)
	}

	// Post quote to Twitter
	if err := client.PostQuote(quote); err != nil {
		log.Fatalf("Failed to post quote: %v", err)
	}

	log.Println("Quote successfully posted to Twitter")
}
